package response

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InvalidArgumentError returns an error with gRPC code InvalidArgument and message.
func InvalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// InternalError returns an error with gRPC code Internal and message.
func InternalError(msg string) error {
	return status.Error(codes.Internal, msg)
}

// NotFoundError returns an error with gRPC code NotFound and message.
func NotFoundError(msg string) error {
	return status.Error(codes.NotFound, msg)
}
