package smtp

import (
	"crypto/tls"
	"net"
	"sync/atomic"
)

type SubmissionsServer struct {
	Addr      string
	TLSConfig *tls.Config

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

	addr := s.Addr
	if addr == "" {
		addr = ":submissions"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return nil
}
