package smtp

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"

	"github.com/quokkamail/quokka/smtp/parser"
)

type session struct {
	srv           *Server
	rwc           net.Conn
	txtReader     *textproto.Reader
	rcptTo        []string
	mailFrom      string
	data          []string
	tls           bool
	authenticated bool
}

func (s *session) serve() {
	s.reply(replyReady("<domain>"))

	s.txtReader = textproto.NewReader(bufio.NewReader(s.rwc))

	for {
		cmdAndArgs, err := s.txtReader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Printf("error: %s\n", err)
			return
		}

		command, _, _ := strings.Cut(cmdAndArgs, " ")

		switch strings.ToUpper(command) {
		case "EHLO":
			s.handleEhloCommand()
		case "HELO":
			s.handleHELOCommand()
		case "MAIL":
			s.handleMailCommand(cmdAndArgs)
		case "RCPT":
			s.handleRCPTCommand(cmdAndArgs)
		case "DATA":
			s.handleDATACommand()
		case "QUIT":
			s.handleQUITCommand()
		case "RSET":
			s.handleRSETCommand()
		case "NOOP":
			s.handleNOOPCommand()
		case "STARTTLS":
			s.handleStartTLSCommand()
		case "AUTH":
			s.handleAuthCommand()
		default:
			s.reply(replyCommandUnrecognized())
		}
	}
}

func (s *session) reset() {
	s.rcptTo = make([]string, 0)
	s.mailFrom = ""
	s.data = make([]string, 0)
}

func (s *session) handleMailCommand(cmdAndArgs string) {
	if s.isNotAuthenticatedWhenMandatory() {
		s.reply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom != "" {
		s.reply(replyBadSequence())
		return
	}

	mailCmd, err := parser.NewMailCommand(cmdAndArgs)
	if err != nil {
		s.reply(replySyntaxError())
		return
	}

	s.mailFrom = mailCmd.ReversePath
	s.reply(replyOk())
}

func (s *session) handleRCPTCommand(cmdAndArgs string) {
	if s.isNotAuthenticatedWhenMandatory() {
		s.reply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom == "" {
		s.reply(replyBadSequence())
		return
	}

	recipientCmd, err := parser.NewRecipientCommand(cmdAndArgs)
	if err != nil {
		s.reply(replySyntaxError())
		return
	}

	s.rcptTo = append(s.rcptTo, recipientCmd.ForwardPath)
	s.reply(replyOk())
}

func (s *session) handleEhloCommand() {
	extensions := []string{"AUTH PLAIN", "PIPELINING"}
	if !s.tls {
		extensions = append(extensions, "STARTTLS")
	}

	s.reply(replyEhloOk(extensions))
}

func (s *session) handleHELOCommand() {
	s.reply(replyHeloOk())
}

func (s *session) handleQUITCommand() {
	s.reply(replyClosingConnection("<domain>"))
	s.rwc.Close()
}

func (s *session) handleRSETCommand() {
	s.reset()
	s.reply(replyOk())
}

func (s *session) handleNOOPCommand() {
	s.reply(replyOk())
}

func (s *session) handleDATACommand() {
	if s.isNotAuthenticatedWhenMandatory() {
		s.reply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom == "" || len(s.rcptTo) == 0 {
		s.reply(replyBadSequence())
		return
	}

	s.reply(replyStartMailInput())

	dl, err := s.txtReader.ReadDotLines()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		log.Printf("error: %s\n", err)
		return
	}

	s.data = dl
	s.reply(replyOk())
}

func (s *session) handleStartTLSCommand() {
	if s.tls {
		s.reply(replyBadSequence())
		return
	}

	s.reply(replyReadyToStartTLS())

	tlsConn := tls.Server(s.rwc, s.srv.TLSConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("error: %s\n", err)
		s.reply(replyTLSNotAvailable())
		return
	}

	s.rwc = tlsConn
	s.txtReader = textproto.NewReader(bufio.NewReader(s.rwc))
	s.tls = true
	s.reset()
}

func (s *session) handleAuthCommand() {
	if s.srv.AuthenticationEncrypted && !s.tls {
		s.reply(replyMustIssueSTARTTLSFirst())
		return
	}
}

func (s *session) isNotAuthenticatedWhenMandatory() bool {
	return s.srv.AuthenticationMandatory && !s.authenticated
}
