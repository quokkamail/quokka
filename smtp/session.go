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
	authenticated bool
	rwc           net.Conn
	srv           *Server
	tls           bool
	txtReader     *textproto.Reader

	data     []string
	mailFrom string
	rcptTo   []string
}

func (s *session) serve() {
	s.reply(replyReady(s.srv.Domain))

	s.txtReader = textproto.NewReader(bufio.NewReader(s.rwc))

	for {
		// s.rwc.SetReadDeadline()

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
			s.handleHeloCommand()
		case "MAIL":
			s.handleMailCommand(cmdAndArgs)
		case "RCPT":
			s.handleRcptCommand(cmdAndArgs)
		case "DATA":
			s.handleDataCommand()
		case "QUIT":
			s.handleQuitCommand()
		case "RSET":
			s.handleRsetCommand()
		case "NOOP":
			s.handleNoopCommand()
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

func (s *session) handleRcptCommand(cmdAndArgs string) {
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

func (s *session) handleHeloCommand() {
	s.reply(replyHeloOk())
}

func (s *session) handleQuitCommand() {
	s.reply(replyClosingConnection(s.srv.Domain))
	s.rwc.Close()
}

func (s *session) handleRsetCommand() {
	s.reset()
	s.reply(replyOk())
}

func (s *session) handleNoopCommand() {
	s.reply(replyOk())
}

func (s *session) handleDataCommand() {
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

	if s.authenticated {
		s.reply(replyBadSequence())
		return
	}
}

func (s *session) isNotAuthenticatedWhenMandatory() bool {
	return s.srv.AuthenticationMandatory && !s.authenticated
}
