package smtp

import (
	"crypto/tls"
	"net"
	"sync/atomic"
)

type SubmissionServer struct {
	Addr      string
	TLSConfig *tls.Config

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

	addr := s.Addr
	if addr == "" {
		addr = ":submission"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return nil
}
