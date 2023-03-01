package smtp

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

var (
	ErrServerClosed           = errors.New("smtp: Server closed")
	ErrMissingServerAddr      = errors.New("smtp: Missing Server Addr")
	ErrMissingServerTLSConfig = errors.New("smtp: Missing Server TLSConfig")
)

// A Server defines parameters for running an SMTP server.
type Server struct {
	Addr      string
	TLSConfig *tls.Config

	AuthenticationEncrypted bool
	AuthenticationMandatory bool

	inShutdown atomic.Bool
	sessions   map[*session]struct{}
	mu         sync.Mutex
}

func (srv *Server) Close() error {
	srv.inShutdown.Store(true)

	for s := range srv.sessions {
		s.rwc.Close()
		delete(srv.sessions, s)
	}

	return nil
}

func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}

	if srv.Addr == "" {
		return ErrMissingServerAddr
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		rw, err := l.Accept()
		if err != nil {
			if srv.shuttingDown() {
				return ErrServerClosed
			}

			return err
		}

		isTLS := false
		if _, ok := rw.(*tls.Conn); ok {
			isTLS = true
		}

		s := &session{srv: srv, rwc: rw, tls: isTLS}
		srv.trackSession(s, true)
		go s.serve()
	}
}

func (srv *Server) ListenAndServeTLS() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}

	if srv.Addr == "" {
		return ErrMissingServerAddr
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return srv.ServeTLS(ln)
}

func (srv *Server) ServeTLS(l net.Listener) error {
	if srv.TLSConfig == nil {
		return ErrMissingServerTLSConfig
	}

	tlsListener := tls.NewListener(l, srv.TLSConfig)
	return srv.Serve(tlsListener)
}

func (srv *Server) Shutdown(ctx context.Context) error {
	srv.inShutdown.Store(true)

	return nil
}

func (srv *Server) shuttingDown() bool {
	return srv.inShutdown.Load()
}

func (srv *Server) trackSession(s *session, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.sessions == nil {
		srv.sessions = make(map[*session]struct{})
	}

	if add {
		srv.sessions[s] = struct{}{}
	} else {
		delete(srv.sessions, s)
	}
}
