package metrics

import "github.com/prometheus/client_golang/prometheus"

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func Path(path string) ServerOption {
	return func(s *Server) {
		s.path = path
	}
}

func Gatherer(gatherer prometheus.Gatherer) ServerOption {
	return func(s *Server) {
		s.gatherer = gatherer
	}
}
