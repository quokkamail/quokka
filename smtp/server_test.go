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

package smtp_test

import (
	"fmt"
	"net"
	"net/textproto"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/quokkamail/quokka/smtp"
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
	srv := &smtp.Server{}

	ls, err := newLocalListener()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.Serve(ls)
		must.NoError(t, err)
	}()

	t.Cleanup(func() { srv.Close() })

	return &TestServer{
		t: t,
		s: srv,
		l: ls,
	}
}
