package smtp

import (
	"net"
	"sync/atomic"
)

type SubmissionServer struct {
	Config Config

	inShutdown atomic.Bool
}

func (s *SubmissionServer) Close() error {
	s.inShutdown.Store(true)

	return nil
}

func (s *SubmissionServer) ListenAndServe() error {
	if s.inShutdown.Load() {
		return ErrServerClosed
	}

	addr := s.Config.Addr
	if addr == "" {
		addr = ":submission"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return s.Serve(ln)
}

func (s *SubmissionServer) Serve(l net.Listener) error {
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}

		c := s.newConn(rw)
		go c.serve()
	}
}

func (s *SubmissionServer) newConn(rwc net.Conn) *conn {
	return &conn{
		config: s.Config,
		conn:   rwc,
	}
}
