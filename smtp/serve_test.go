// Copyright 2023 Quokka Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// End-to-end tests

package smtp_test

import (
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"testing"

	"github.com/quokkamail/quokka/smtp"
	"github.com/shoenig/test/must"
)

func newLocalListener() (net.Listener, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			return nil, fmt.Errorf("failed to listen on a port: %w", err)
		}
	}

	return l, nil
}

type TestServer struct {
	t testing.TB
	s *smtp.Server
	l net.Listener
}

func (ts *TestServer) Client() *textproto.Conn {
	conn, err := net.Dial("tcp", ts.l.Addr().String())
	if err != nil {
		ts.t.Fatalf("failed to dial: %v", err)
	}

	textConn := textproto.NewConn(conn)

	ts.t.Cleanup(func() {
		textConn.Close()
	})

	return textConn
}

func NewTestServer(t testing.TB) *TestServer {
	srv := &smtp.Server{
		Domain: "quokka.test",
	}

	ls, err := newLocalListener()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.Serve(ls)
		must.NoError(t, err)
	}()

	t.Cleanup(func() {
		srv.Close()
	})

	return &TestServer{
		t: t,
		s: srv,
		l: ls,
	}
}

func TestCommandUnrecognized(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	_, _, err := tc.ReadResponse(220)
	must.NoError(t, err)

	_, err = tc.Cmd("DUMMYCOMMAND")
	must.NoError(t, err)

	_, message, err := tc.ReadResponse(500)
	must.NoError(t, err)
	must.Eq(t, "5.5.2 Syntax error, command unrecognized", message)
}

func TestHeloCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	_, _, err := tc.ReadResponse(220)
	must.NoError(t, err)

	_, err = tc.Cmd("HELO")
	must.NoError(t, err)

	_, message, err := tc.ReadResponse(250)
	must.NoError(t, err)
	must.Eq(t, "Hello, nice to meet you", message)
}

func TestQuitCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	_, _, err := tc.ReadResponse(220)
	must.NoError(t, err)

	_, err = tc.Cmd("QUIT")
	must.NoError(t, err)

	_, message, err := tc.ReadResponse(221)
	must.NoError(t, err)
	must.Eq(t, "2.0.0 quokka.test service closing transmission channel", message)
}

func TestRsetCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	_, _, err := tc.ReadResponse(220)
	must.NoError(t, err)

	_, err = tc.Cmd("RSET")
	must.NoError(t, err)

	_, message, err := tc.ReadResponse(250)
	must.NoError(t, err)
	must.Eq(t, "2.0.0 Requested mail action okay, completed", message)
}

func TestNoopCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	_, _, err := tc.ReadResponse(220)
	must.NoError(t, err)

	_, err = tc.Cmd("NOOP")
	must.NoError(t, err)

	_, message, err := tc.ReadResponse(250)
	must.NoError(t, err)
	must.Eq(t, "2.0.0 Requested mail action okay, completed", message)
}

func TestMailCommand(t *testing.T) {
	t.Parallel()

	type command struct {
		command     string
		wantCode    int
		wantMessage string
	}

	testCases := []struct {
		name     string
		commands []command
	}{
		{
			name: "Valid",
			commands: []command{
				{
					command:     "MAIL FROM: <mail@domain.ext>",
					wantCode:    250,
					wantMessage: "2.1.0 Requested mail action okay, completed",
				},
			},
		},
		{
			name: "Invalid",
			commands: []command{
				{
					command:     "MAIL",
					wantCode:    501,
					wantMessage: "5.5.4 Syntax error in parameters or arguments",
				},
			},
		},
		{
			name: "BadSequence",
			commands: []command{
				{
					command:     "MAIL FROM: <mail@domain.ext>",
					wantCode:    250,
					wantMessage: "2.1.0 Requested mail action okay, completed",
				},
				{
					command:     "MAIL FROM: <mail@domain.ext>",
					wantCode:    503,
					wantMessage: "5.5.1 Bad sequence of commands",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewTestServer(t)
			c := s.Client()

			_, _, err := c.ReadResponse(220)
			must.NoError(t, err)

			for _, cmd := range tc.commands {
				_, err := c.Cmd(cmd.command)
				must.NoError(t, err)

				_, message, err := c.ReadResponse(cmd.wantCode)
				must.NoError(t, err)
				must.Eq(t, cmd.wantMessage, message)
			}
		})
	}
}

func TestRcptCommand(t *testing.T) {
	t.Parallel()

	type command struct {
		command     string
		wantCode    int
		wantMessage string
	}

	testCases := []struct {
		name     string
		commands []command
	}{
		{
			name: "Valid",
			commands: []command{
				{
					command:     "MAIL FROM: <mail@domain.ext>",
					wantCode:    250,
					wantMessage: "2.1.0 Requested mail action okay, completed",
				},
				{
					command:     "RCPT TO: <mail@domain.ext>",
					wantCode:    250,
					wantMessage: "2.1.5 Requested mail action okay, completed",
				},
			},
		},
		{
			name: "Invalid",
			commands: []command{
				{
					command:     "MAIL FROM: <mail@domain.ext>",
					wantCode:    250,
					wantMessage: "2.1.0 Requested mail action okay, completed",
				},
				{
					command:     "RCPT",
					wantCode:    501,
					wantMessage: "5.5.4 Syntax error in parameters or arguments",
				},
			},
		},
		{
			name: "BadSequence",
			commands: []command{
				{
					command:     "RCPT TO: <mail@domain.ext>",
					wantCode:    503,
					wantMessage: "5.5.1 Bad sequence of commands",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewTestServer(t)
			c := s.Client()

			_, _, err := c.ReadResponse(220)
			must.NoError(t, err)

			for _, cmd := range tc.commands {
				_, err := c.Cmd(cmd.command)
				must.NoError(t, err)

				_, message, err := c.ReadResponse(cmd.wantCode)
				must.NoError(t, err)
				must.Eq(t, cmd.wantMessage, message)
			}
		})
	}
}

func TestEhloCommand(t *testing.T) {
	t.Parallel()

	wantMessages := []string{
		"Hello, nice to meet you",
		"AUTH PLAIN",
		"ENHANCEDSTATUSCODES",
		"PIPELINING",
		"STARTTLS",
	}

	ts := NewTestServer(t)
	tc := ts.Client()

	_, _, err := tc.ReadResponse(220)
	must.NoError(t, err)

	_, err = tc.Cmd("EHLO")
	must.NoError(t, err)

	_, message, err := tc.ReadResponse(250)
	must.NoError(t, err)
	must.Eq(t, strings.Join(wantMessages, "\n"), message)
}
