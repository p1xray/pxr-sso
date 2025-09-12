package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

const (
	defaultPort = "80"
	network     = "tcp"
)

// Server provides access to the gRPC server.
type Server struct {
	App     *grpc.Server
	notify  chan error
	address string
}

// New returns new gRPC server instance.
func New(opts ...Option) *Server {
	s := &Server{
		App:     grpc.NewServer(),
		notify:  make(chan error),
		address: net.JoinHostPort("", defaultPort),
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start - starts the gRPC server.
func (s *Server) Start() {
	go func() {
		defer close(s.notify)

		ln, err := net.Listen(network, s.address)
		if err != nil {
			s.notify <- fmt.Errorf("failed to listen: %w", err)
			return
		}

		s.notify <- s.App.Serve(ln)
	}()
}

// Notify - notifies about gRPC server errors.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Stop - stops the gRPC server.
func (s *Server) Stop() {
	s.App.GracefulStop()
}
