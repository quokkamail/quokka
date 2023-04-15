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
