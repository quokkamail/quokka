package smtp

import (
	"crypto/tls"
	"errors"
	"net"
	"sync/atomic"
)

var (
	ErrServerClosed = errors.New("smtp: Server closed")
)

type Config struct {
	Addr      string
	TLSConfig *tls.Config

	AuthenticationEncrypted bool
	AuthenticationMandatory bool
}

type Server struct {
	Config

	inShutdown atomic.Bool
}

func (s *Server) Close() error {
	s.inShutdown.Store(true)

	return nil
}

func (s *Server) ListenAndServe() error {
	if s.inShutdown.Load() {
		return ErrServerClosed
	}

	addr := s.Config.Addr
	if addr == "" {
		addr = ":smtp"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return s.Serve(ln)
}

func (s *Server) Serve(l net.Listener) error {
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}

		c := s.newConn(rw)
		go c.serve()
	}
}

func (s *Server) newConn(rwc net.Conn) *session {
	return &session{
		config: s.Config,
		conn:   rwc,
	}
}
