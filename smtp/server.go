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

package smtp

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrMissingServerAddr      = errors.New("smtp: Missing Server Addr")
	ErrMissingServerTLSConfig = errors.New("smtp: Missing Server TLSConfig")
	ErrServerClosed           = errors.New("smtp: Server closed")
)

// A Server defines parameters for running an SMTP server.
type Server struct {
	Address   string
	TLSConfig *tls.Config

	AuthenticationEncrypted bool
	AuthenticationMandatory bool
	Domain                  string

	Initial220MessageTimeout time.Duration
	MailCommandTimeout       time.Duration
	RcptCommandTimeout       time.Duration

	inShutdown atomic.Bool
	mu         sync.Mutex
	sessions   map[*session]struct{}
}

func (srv *Server) Close() error {
	srv.inShutdown.Store(true)

	srv.mu.Lock()
	defer srv.mu.Unlock()

	for s := range srv.sessions {
		s.conn.Close()
		delete(srv.sessions, s)
	}

	return nil
}

func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}

	if srv.Address == "" {
		return ErrMissingServerAddr
	}

	ln, err := net.Listen("tcp", srv.Address)
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
			// if srv.shuttingDown() {
			// 	return ErrServerClosed
			// }

			return err
		}

		_, isTLS := rw.(*tls.Conn)

		s := &session{srv: srv, conn: rw, tls: isTLS}
		srv.trackSession(s, true)
		go s.serve()
	}
}

func (srv *Server) ListenAndServeTLS() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}

	if srv.Address == "" {
		return ErrMissingServerAddr
	}

	ln, err := net.Listen("tcp", srv.Address)
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

// func (srv *Server) Shutdown(ctx context.Context) error {
// 	srv.inShutdown.Store(true)

// 	return nil
// }

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
