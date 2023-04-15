package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"
)

type Server struct {
	Address string
	srv     *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}

	return s.srv.Shutdown(ctx)
}

func (s *Server) ListenAndServe() error {
	mux := &http.ServeMux{}
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.srv = &http.Server{
		Addr:              s.Address,
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Minute,
	}

	return s.srv.ListenAndServe()
}
