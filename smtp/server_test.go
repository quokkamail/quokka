package smtp_test

import (
	"fmt"
	"net"
	"net/textproto"
	"testing"

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

type clientServerTest struct {
	// t testing.TB
	s *smtp.Server
	c *textproto.Conn
}

func newClientServerTest(t testing.TB) *clientServerTest {
	cst := &clientServerTest{
		// t: t,
		s: &smtp.Server{},
	}

	ls, err := newLocalListener()
	if err != nil {
		t.Fatal(err)
	}

	go cst.s.Serve(ls)

	conn, err := net.Dial("tcp", ls.Addr().String())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	cst.c = textproto.NewConn(conn)

	t.Cleanup(func() {
		cst.c.Close()
		cst.s.Close()
	})

	return cst
}

type TestServer struct {
	t testing.TB
	s *smtp.Server
	l net.Listener
	// c *textproto.Conn
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

	go srv.Serve(ls)

	t.Cleanup(func() {
		srv.Close()
	})

	return &TestServer{
		t: t,
		s: srv,
		l: ls,
	}
}
