package grpcserver

import "net"

// Option is how options for the Server are set up.
type Option func(*Server)

// WithPort sets up a port for gRPC server.
func WithPort(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}
