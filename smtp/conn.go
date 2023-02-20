package smtp

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"

	"github.com/quokkamail/quokka/smtp/parser"
)

type session struct {
	config        Config
	conn          net.Conn
	txtReader     *textproto.Reader
	rcptTo        []string
	mailFrom      string
	data          []string
	tls           bool
	authenticated bool
}

func (s *session) serve() {
	s.replyWithCode(220)

	s.txtReader = textproto.NewReader(bufio.NewReader(s.conn))

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
			s.handleEHLOCommand()

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
			s.replyWithCode(500)
		}
	}
}

func (s *session) reset() {
	s.rcptTo = make([]string, 0)
	s.mailFrom = ""
	s.data = make([]string, 0)
}

func (s *session) replyWithCode(code uint) {
	var message string

	switch code {
	case 220:
		message = "<domain> Service ready"
	case 221:
		message = "<domain> Service closing transmission channel"
	case 250:
		message = "Requested mail action okay, completed"

	case 354:
		message = "Start mail input; end with <CRLF>.<CRLF>"

	case 454:
		message = "TLS not available due to temporary reason"

	case 500:
		message = "Syntax error, command unrecognized (This may include errors such as command line too long)"
	case 501:
		message = "Syntax error in parameters or arguments"
	case 502:
		message = "Command not implemented"
	case 503:
		message = "Bad sequence of commands"
	case 504:
		message = "Command parameter not implemented"
	case 530:
		message = "Must issue a STARTTLS / AUTH command first"
	}

	s.replyWithCodeAndMessage(code, message)
}

func (s *session) replyWithCodeAndMessage(code uint, message string) {
	fmt.Fprintf(s.conn, "%d %s\r\n", code, message)
}

func (s *session) replyWithCodeAndMessages(code uint, messages []string) {
	for i, m := range messages {
		sep := "-"
		if i == len(messages)-1 {
			sep = " "
		}

		fmt.Fprintf(s.conn, "%d%s%s\r\n", code, sep, m)
	}
}

func (s *session) handleMailCommand(cmdAndArgs string) {
	if s.config.AuthenticationMandatory && !s.authenticated {
		s.replyWithCode(530)
		return
	}

	if s.mailFrom != "" {
		s.replyWithCode(503)
		return
	}

	mailCmd, err := parser.NewMailCommand(cmdAndArgs)
	if err != nil {
		s.replyWithCode(501)
		return
	}

	s.mailFrom = mailCmd.ReversePath
	s.replyWithCode(250)
}

func (s *session) handleRCPTCommand(cmdAndArgs string) {
	if s.config.AuthenticationMandatory && !s.authenticated {
		s.replyWithCode(530)
		return
	}

	if s.mailFrom == "" {
		s.replyWithCode(503)
		return
	}

	recipientCmd, err := parser.NewRecipientCommand(cmdAndArgs)
	if err != nil {
		s.replyWithCode(501)
		return
	}

	s.rcptTo = append(s.rcptTo, recipientCmd.ForwardPath)
	s.replyWithCode(250)
}

func (s *session) handleEHLOCommand() {
	msgs := []string{
		"Hello, nice to meet you",
	}

	msgs = append(msgs, "AUTH PLAIN")
	msgs = append(msgs, "PIPELINING")

	if !s.tls {
		msgs = append(msgs, "STARTTLS")
	}

	s.replyWithCodeAndMessages(250, msgs)
}

func (s *session) handleHELOCommand() {
	s.replyWithCodeAndMessage(250, "Hello, nice to meet you")
}

func (s *session) handleQUITCommand() {
	s.replyWithCode(221)
	s.conn.Close()
}

func (s *session) handleRSETCommand() {
	s.reset()
	s.replyWithCode(250)
}

func (s *session) handleNOOPCommand() {
	s.replyWithCode(250)
}

func (s *session) handleDATACommand() {
	if s.config.AuthenticationMandatory && !s.authenticated {
		s.replyWithCode(530)
		return
	}

	if s.mailFrom == "" || len(s.rcptTo) == 0 {
		s.replyWithCode(503)
		return
	}

	s.replyWithCode(354)

	dl, err := s.txtReader.ReadDotLines()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		log.Printf("error: %s\n", err)
		return
	}

	s.data = dl
	s.replyWithCode(250)
}

func (s *session) handleStartTLSCommand() {
	if s.tls {
		s.replyWithCode(503)
		return
	}

	s.replyWithCodeAndMessage(220, "Ready to start TLS")

	tlsConn := tls.Server(s.conn, s.config.TLSConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("error: %s\n", err)
		s.replyWithCode(454)
		return
	}

	s.conn = tlsConn
	s.txtReader = textproto.NewReader(bufio.NewReader(s.conn))
	s.tls = true
	s.reset()
}

func (s *session) handleAuthCommand() {
	if s.config.AuthenticationEncrypted && !s.tls {
		s.replyWithCode(530)
		return
	}
}
