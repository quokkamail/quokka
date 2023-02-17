package smtp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
)

const (
	cmdDATA string = "DATA"
	cmdEHLO string = "EHLO"
	cmdHELO string = "HELO"
	cmdMAIL string = "MAIL"
	cmdQUIT string = "QUIT"
	cmdRCPT string = "RCPT"
	cmdRSET string = "RSET"
	cmdNOOP string = "NOOP"
)

type replyCode int

const (
	replyCode220 replyCode = 220
	replyCode221 replyCode = 221
	replyCode250 replyCode = 250

	replyCode354 replyCode = 354

	replyCode500 replyCode = 500
	replyCode501 replyCode = 501
	replyCode502 replyCode = 502
	replyCode503 replyCode = 503
	replyCode504 replyCode = 504
)

type conn struct {
	server   *Server
	rwc      net.Conn
	rcptTo   []string
	mailFrom []string
	data     []string
}

func (c *conn) serve() {
	c.replyWithCode(replyCode220)

	r := textproto.NewReader(bufio.NewReader(c.rwc))

	for {
		l, err := r.ReadLine()
		if err != nil {
			// if errors.Is(err, io.EOF) {
			// 	return
			// }

			log.Printf("error: %s\n", err)
			return
		}

		l = strings.ToUpper(l)

		cmd, arg, hasArg := strings.Cut(l, " ")

		switch cmd {
		case cmdEHLO:
			c.handleEHLOCommand()

		case cmdHELO:
			c.handleHELOCommand()

		case cmdMAIL:
			if !hasArg {
				c.replyWithCode(replyCode501)
				continue
			}

			c.handleMAILCommand(arg)

		case cmdRCPT:
			if !hasArg {
				c.replyWithCode(replyCode501)
				continue
			}

			c.handleRCPTCommand(arg)

		case cmdDATA:
			c.handleDATACommand(r)

		case cmdQUIT:
			c.handleQUITCommand()

		case cmdRSET:
			c.handleRSETCommand()

		case cmdNOOP:
			c.handleNOOPCommand()

		default:
			c.replyWithCode(replyCode500)
		}
	}
}

func (c *conn) reset() {
	c.rcptTo = make([]string, 0)
	c.mailFrom = make([]string, 0)
	c.data = make([]string, 0)
}

func (c *conn) replyWithCode(code replyCode) {
	var message string

	switch code {
	case replyCode220:
		message = "<domain> Service ready"
	case replyCode221:
		message = "<domain> Service closing transmission channel"
	case replyCode250:
		message = "Requested mail action okay, completed"

	case replyCode354:
		message = "Start mail input; end with <CRLF>.<CRLF>"

	case replyCode500:
		message = "Syntax error, command unrecognized (This may include errors such as command line too long)"
	case replyCode501:
		message = "Syntax error in parameters or arguments"
	case replyCode502:
		message = "Command not implemented"
	case replyCode503:
		message = "Bad sequence of commands"
	case replyCode504:
		message = "Command parameter not implemented"
	}

	c.replyWithCodeAndMessage(code, message)
}

func (c *conn) replyWithCodeAndMessage(code replyCode, message string) {
	fmt.Fprintf(c.rwc, "%d %s\r\n", code, message)
}

func (c *conn) handleMAILCommand(arg string) {
	if len(c.mailFrom) > 0 {
		c.replyWithCode(replyCode503)
		return
	}

	if len(arg) < 5 || arg[:5] != "FROM:" {
		c.replyWithCode(replyCode501)
		return
	}

	mailFrom := strings.Split(strings.Trim(arg[5:], " "), " ")
	if mailFrom[0] == "" {
		c.replyWithCode(replyCode501)
	}

	c.mailFrom = mailFrom
	c.replyWithCode(replyCode250)
}

func (c *conn) handleRCPTCommand(arg string) {
	if len(c.mailFrom) == 0 {
		c.replyWithCode(replyCode503)
		return
	}

	if len(arg) < 3 || arg[:3] != "TO:" {
		c.replyWithCode(replyCode501)
		return
	}

	rcptTo := strings.Split(strings.Trim(arg[3:], " "), " ")
	if rcptTo[0] == "" {
		c.replyWithCode(replyCode501)
	}

	c.rcptTo = append(c.rcptTo, rcptTo...)
	c.replyWithCode(replyCode250)
}

func (c *conn) handleEHLOCommand() {
	c.replyWithCodeAndMessage(replyCode250, "Hello, nice to meet you")
}

func (c *conn) handleHELOCommand() {
	c.replyWithCodeAndMessage(replyCode250, "Hello, nice to meet you")
}

func (c *conn) handleQUITCommand() {
	c.replyWithCode(replyCode221)
	// c.rwc.Close()
}

func (c *conn) handleRSETCommand() {
	c.reset()
	c.replyWithCode(replyCode250)
}

func (c *conn) handleNOOPCommand() {
	c.replyWithCode(replyCode250)
}

func (c *conn) handleDATACommand(r *textproto.Reader) {
	if len(c.mailFrom) == 0 || len(c.rcptTo) == 0 {
		c.replyWithCode(replyCode503)
		return
	}

	c.replyWithCode(replyCode354)

	dl, err := r.ReadDotLines()
	if err != nil {
		// if errors.Is(err, io.EOF) {
		// 	return
		// }

		log.Printf("error: %s\n", err)
		return
	}

	c.data = dl
	c.replyWithCode(replyCode250)
}
