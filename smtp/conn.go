package smtp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
)

type conn struct {
	server   *Server
	rwc      net.Conn
	rcptTo   []string
	mailFrom string
	data     []string
}

func (c *conn) serve() {
	c.replyWithCode(220)

	r := textproto.NewReader(bufio.NewReader(c.rwc))

	for {
		cmdAndArgs, err := r.ReadLine()
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
			// if !hasArg {
			// 	c.replyWithCode(replyCode501)
			// 	continue
			// }

			c.handleRCPTCommand(cmdAndArgs)

		case "DATA":
			c.handleDATACommand(r)

		case "QUIT":
			c.handleQUITCommand()

		case "RSET":
			c.handleRSETCommand()

		case "NOOP":
			c.handleNOOPCommand()

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
	fmt.Fprintf(c.rwc, "%d %s\r\n", code, message)
}

func (c *conn) handleMailCommand(cmdAndArgs string) {
	if c.mailFrom != "" {
		c.replyWithCode(503)
		return
	}

	mailCmd, err := ParseMailCommand(cmdAndArgs)
	if err != nil {
		c.replyWithCode(501)
		return
	}

	c.mailFrom = mailCmd.ReversePath
	c.replyWithCode(250)
}

func (c *conn) handleRCPTCommand(cmdAndArgs string) {
	if len(c.mailFrom) == 0 {
		c.replyWithCode(503)
		return
	}

	if len(cmdAndArgs) < 3 || cmdAndArgs[:3] != "TO:" {
		c.replyWithCode(501)
		return
	}

	rcptTo := strings.Split(strings.Trim(cmdAndArgs[3:], " "), " ")
	if rcptTo[0] == "" {
		c.replyWithCode(501)
	}

	c.rcptTo = append(c.rcptTo, rcptTo...)
	c.replyWithCode(250)
}

func (c *conn) handleEHLOCommand() {
	c.replyWithCodeAndMessage(250, "Hello, nice to meet you")
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

func (c *conn) handleDATACommand(r *textproto.Reader) {
	if len(c.mailFrom) == 0 || len(c.rcptTo) == 0 {
		c.replyWithCode(503)
		return
	}

	c.replyWithCode(354)

	dl, err := r.ReadDotLines()
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
