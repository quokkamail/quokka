package smtp

import (
	"crypto/tls"
	"net"
	"sync/atomic"
)

type SubmissionsServer struct {
	Config Config

	inShutdown atomic.Bool
}

func (s *SubmissionsServer) Close() error {
	s.inShutdown.Store(true)

	return nil
}

func (s *SubmissionsServer) ListenAndServeTLS() error {
	if s.inShutdown.Load() {
		return ErrServerClosed
	}

	addr := s.Config.Addr
	if addr == "" {
		addr = ":465"
	}

	ln, err := tls.Listen("tcp", addr, s.Config.TLSConfig)
	if err != nil {
		return err
	}

	defer ln.Close()

	return s.Serve(ln)
}

func (s *SubmissionsServer) Serve(l net.Listener) error {
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}

		c := s.newConn(rw)
		go c.serve()
	}
}

func (s *SubmissionsServer) newConn(rwc net.Conn) *conn {
	return &conn{
		conn: rwc,
		tls:  true,
	}
}
