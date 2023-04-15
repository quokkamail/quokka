package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	promauto.NewCounter(prometheus.CounterOpts{
		Name: "quokka_test_total",
		Help: "Quokka test metric",
	})

	mux := &http.ServeMux{}
	mux.Handle("/metrics", promhttp.Handler())

	s.srv = &http.Server{
		Addr:              s.Address,
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Minute,
	}

	return s.srv.ListenAndServe()
}
