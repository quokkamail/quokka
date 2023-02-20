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

type conn struct {
	server    *Server
	conn      net.Conn
	txtReader *textproto.Reader
	rcptTo    []string
	mailFrom  string
	data      []string
	tls       bool
}

func (c *conn) serve() {
	c.replyWithCode(220)

	c.txtReader = textproto.NewReader(bufio.NewReader(c.conn))

	for {
		cmdAndArgs, err := c.txtReader.ReadLine()
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
			c.handleEHLOCommand()

		case "HELO":
			c.handleHELOCommand()

		case "MAIL":
			c.handleMailCommand(cmdAndArgs)

		case "RCPT":
			c.handleRCPTCommand(cmdAndArgs)

		case "DATA":
			c.handleDATACommand()

		case "QUIT":
			c.handleQUITCommand()

		case "RSET":
			c.handleRSETCommand()

		case "NOOP":
			c.handleNOOPCommand()

		case "STARTTLS":
			c.handleStartTLSCommand()

		default:
			c.replyWithCode(500)
		}
	}
}

func (c *conn) reset() {
	c.rcptTo = make([]string, 0)
	c.mailFrom = ""
	c.data = make([]string, 0)
}

func (c *conn) replyWithCode(code uint) {
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
	}

	c.replyWithCodeAndMessage(code, message)
}

func (c *conn) replyWithCodeAndMessage(code uint, message string) {
	fmt.Fprintf(c.conn, "%d %s\r\n", code, message)
}

func (c *conn) replyWithCodeAndMessages(code uint, messages []string) {
	for i, m := range messages {
		sep := "-"
		if i == len(messages)-1 {
			sep = " "
		}

		fmt.Fprintf(c.conn, "%d%s%s\r\n", code, sep, m)
	}
}

func (c *conn) handleMailCommand(cmdAndArgs string) {
	if c.mailFrom != "" {
		c.replyWithCode(503)
		return
	}

	mailCmd, err := parser.NewMailCommand(cmdAndArgs)
	if err != nil {
		c.replyWithCode(501)
		return
	}

	c.mailFrom = mailCmd.ReversePath
	c.replyWithCode(250)
}

func (c *conn) handleRCPTCommand(cmdAndArgs string) {
	if c.mailFrom == "" {
		c.replyWithCode(503)
		return
	}

	recipientCmd, err := parser.NewRecipientCommand(cmdAndArgs)
	if err != nil {
		c.replyWithCode(501)
		return
	}

	c.rcptTo = append(c.rcptTo, recipientCmd.ForwardPath)
	c.replyWithCode(250)
}

func (c *conn) handleEHLOCommand() {
	msgs := []string{
		"Hello, nice to meet you",
	}

	if !c.tls {
		msgs = append(msgs, "STARTTLS")
	}

	c.replyWithCodeAndMessages(250, msgs)
}

func (c *conn) handleHELOCommand() {
	c.replyWithCodeAndMessage(250, "Hello, nice to meet you")
}

func (c *conn) handleQUITCommand() {
	c.replyWithCode(221)
	// c.rwc.Close()
}

func (c *conn) handleRSETCommand() {
	c.reset()
	c.replyWithCode(250)
}

func (c *conn) handleNOOPCommand() {
	c.replyWithCode(250)
}

func (c *conn) handleDATACommand() {
	if c.mailFrom == "" || len(c.rcptTo) == 0 {
		c.replyWithCode(503)
		return
	}

	c.replyWithCode(354)

	dl, err := c.txtReader.ReadDotLines()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		log.Printf("error: %s\n", err)
		return
	}

	c.data = dl
	c.replyWithCode(250)
}

func (c *conn) handleStartTLSCommand() {
	c.replyWithCodeAndMessage(220, "Ready to start TLS")

	tlsConn := tls.Server(c.conn, c.server.TLSConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("error: %s\n", err)
		c.replyWithCode(454)
		return
	}

	c.conn = tlsConn
	c.txtReader = textproto.NewReader(bufio.NewReader(c.conn))
	c.tls = true
	c.reset()
}
