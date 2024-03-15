package metrics

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	*http.Server
	gatherer prometheus.Gatherer
	address  string
	path     string
	router   *mux.Router
}

type MetricsServer = Server

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address:  ":0",
		path:     "/metrics",
		gatherer: prometheus.DefaultGatherer,
	}
	for _, o := range opts {
		o(srv)
	}
	srv.router = mux.NewRouter()

	registerMetricsRoute(srv.router, srv.path, srv.gatherer)

	srv.Server = &http.Server{
		Addr:              srv.address,
		ReadHeaderTimeout: 1 * time.Second,
		Handler:           srv.router,
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func registerMetricsRoute(router *mux.Router, path string, gatherer prometheus.Gatherer) {
	router.Handle(path, promhttp.HandlerFor(
		gatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}))
}
