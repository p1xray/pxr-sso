package grpcserver

import "net"

// Option is how options for the Server are set up.
type Option func(*server)

// WithPort sets up a port for gRPC server.
func WithPort(port string) Option {
	return func(s *server) {
		s.address = net.JoinHostPort("", port)
	}
}
