package smtp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	replyCode500 replyCode = 500
	replyCode501 replyCode = 501
	replyCode502 replyCode = 502
	replyCode503 replyCode = 503
	replyCode504 replyCode = 504
)

type conn struct {
	server          *Server
	rwc             net.Conn
	isReceivingData bool
	rcptTo          string
	mailFrom        string
	greeted         bool
	data            []string
}

func (c *conn) serve() {
	// log.Println("New connection")

	// greet the client
	c.replyWithCode(replyCode220)

	r := textproto.NewReader(bufio.NewReader(c.rwc))

	for {
		if c.isReceivingData {
			// log.Println("reading data...")

			dl, err := r.ReadDotLines()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				// log.Printf("error: %s\n", err)
				return
			}

			// log.Printf("data: %s\n", dl)

			c.replyWithCode(replyCode250)

			c.data = dl
			c.isReceivingData = false

			continue
		}

		l, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			// log.Printf("error: %s\n", err)
			return
		}

		l = strings.ToUpper(l)
		// log.Printf("command: %s\n", l)

		cmdAndArg := strings.SplitN(l, " ", 2)

		switch cmdAndArg[0] {
		case cmdEHLO:
			fmt.Fprint(c.rwc, "250 Hello, nice to meet you\r\n")

			// dummy extensions
			// fmt.Fprint(c.rwc, "250-8BITMIME\r\n")
			// fmt.Fprint(c.rwc, "250-SIZE\r\n")
			// fmt.Fprint(c.rwc, "250 PIPELINING\r\n")
			// fmt.Fprint(c.rwc, "250 AUTH\r\n")

			c.greeted = true

		case cmdHELO:
			fmt.Fprintf(c.rwc, "250 Hello, nice to meet you\r\n")

			c.greeted = true

		case cmdMAIL:
			if len(cmdAndArg) == 1 {
				c.replyWithCode(replyCode501)
				continue
			}

			c.handleMailCommand(cmdAndArg[1])

		case cmdRCPT:
			if c.mailFrom == "" {
				fmt.Fprint(c.rwc, "503 Must have sender before recipient\r\n")
				continue
			}

			c.replyWithCode(replyCode250)

			c.rcptTo = "dummy"

		case cmdDATA:
			if c.rcptTo == "" || c.mailFrom == "" {
				fmt.Fprint(c.rwc, "503 Must have valid receiver and originator\r\n")
				continue
			}

			fmt.Fprint(c.rwc, "354 Start mail input; end with <CRLF>.<CRLF>\r\n")

			c.isReceivingData = true

		case cmdQUIT:
			c.replyWithCode(replyCode221)
			c.rwc.Close()

		case cmdRSET:
			c.reset()
			c.replyWithCode(replyCode250)

		case cmdNOOP:
			c.replyWithCode(replyCode250)

		default:
			c.replyWithCode(replyCode500)
		}
	}
}

func (c *conn) reset() {
	c.rcptTo = ""
	c.mailFrom = ""
	c.data = make([]string, 0)
}

func (c *conn) replyWithCode(code replyCode) {
	var text string

	switch code {
	case replyCode220:
		text = "<domain> Service ready"
	case replyCode221:
		text = "<domain> Service closing transmission channel"
	case replyCode250:
		text = "Requested mail action okay, completed"

	case replyCode500:
		text = "Syntax error, command unrecognized (This may include errors such as command line too long)"
	case replyCode501:
		text = "Syntax error in parameters or arguments"
	case replyCode502:
		text = "Command not implemented"
	case replyCode503:
		text = "Bad sequence of commands"
	case replyCode504:
		text = "Command parameter not implemented"
	}

	fmt.Fprintf(c.rwc, "%d %s\r\n", code, text)
}

func (c *conn) handleMailCommand(arg string) {
	if len(arg) < 6 || arg[0:5] != "FROM:" {
		c.replyWithCode(replyCode501)
		return
	}

	_ = strings.Split(strings.Trim(arg[5:], " "), " ")
}
