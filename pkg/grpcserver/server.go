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
type Server interface {
	Start()
	Stop()
	Notify() <-chan error
	Registrar() grpc.ServiceRegistrar
}

// server provides access to the gRPC server.
type server struct {
	innerServer *grpc.Server
	notify      chan error
	address     string
}

// New returns new gRPC server instance.
func New(opts ...Option) *server {
	s := &server{
		innerServer: grpc.NewServer(),
		notify:      make(chan error),
		address:     net.JoinHostPort("", defaultPort),
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start - starts the gRPC server.
func (s *server) Start() {
	go func() {
		defer close(s.notify)

		ln, err := net.Listen(network, s.address)
		if err != nil {
			s.notify <- fmt.Errorf("failed to listen: %w", err)
			return
		}

		s.notify <- s.innerServer.Serve(ln)
	}()
}

// Notify - notifies about gRPC server errors.
func (s *server) Notify() <-chan error {
	return s.notify
}

// Stop - stops the gRPC server.
func (s *server) Stop() {
	s.innerServer.GracefulStop()
}

func (s *server) Registrar() grpc.ServiceRegistrar {
	return s.innerServer
}
