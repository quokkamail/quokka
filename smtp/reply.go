package smtp

import "fmt"

type reply struct {
	code  int
	lines []string
}

func replyReady(domain string) reply {
	return reply{code: 220, lines: []string{fmt.Sprintf("%s Service ready", domain)}}
}

func replyReadyToStartTLS() reply {
	return reply{code: 220, lines: []string{"Ready to start TLS"}}
}

func replyClosingConnection(domain string) reply {
	return reply{code: 221, lines: []string{fmt.Sprintf("%s Service closing transmission channel", domain)}}
}

func replyOk() reply {
	return reply{code: 250, lines: []string{"Requested mail action okay, completed"}}
}

func replyHeloOk() reply {
	return reply{code: 250, lines: []string{"Hello, nice to meet you"}}
}

func replyEhloOk(extensions []string) reply {
	lines := []string{"Hello, nice to meet you"}
	lines = append(lines, extensions...)
	return reply{code: 250, lines: lines}
}

func replyStartMailInput() reply {
	return reply{code: 354, lines: []string{"Start mail input; end with <CRLF>.<CRLF>"}}
}

func replyTLSNotAvailable() reply {
	return reply{code: 454, lines: []string{"TLS not available due to temporary reason"}}
}

func replyCommandUnrecognized() reply {
	return reply{code: 500, lines: []string{"Syntax error, command unrecognized (This may include errors such as command line too long)"}}
}

func replySyntaxError() reply {
	return reply{code: 501, lines: []string{"Syntax error in parameters or arguments"}}
}

// nolint:unused
func replyCommandNotImplemented() reply {
	return reply{code: 502, lines: []string{"Command not implemented"}}
}

func replyBadSequence() reply {
	return reply{code: 503, lines: []string{"Bad sequence of commands"}}
}

// nolint:unused
func replyCommandParameterNotImplemented() reply {
	return reply{code: 504, lines: []string{"Command parameter not implemented"}}
}

func replyMustIssueSTARTTLSFirst() reply {
	return reply{code: 530, lines: []string{"Must issue a STARTTLS command first"}}
}

func replyAuthenticationRequired() reply {
	return reply{code: 530, lines: []string{"Authentication required"}}
}

func replyAuthenticationSucceeded() reply {
	return reply{code: 235, lines: []string{"Authentication Succeeded"}}
}

func (s *session) reply(r reply) {
	for _, m := range r.lines[:len(r.lines)-1] {
		fmt.Fprintf(s.rwc, "%d-%s\r\n", r.code, m)
	}
	fmt.Fprintf(s.rwc, "%d %s\r\n", r.code, r.lines[len(r.lines)-1])
}
