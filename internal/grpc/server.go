package ssoserver

import (
	"context"
	ssopb "github.com/p1xray/pxr-sso-protos/gen/go/sso"
	"google.golang.org/grpc"
	"pxr-sso/internal/service"
)

type serverAPI struct {
	ssopb.UnimplementedSsoServer
	auth service.AuthService
}

// Register registers the implementation of the API service with the gRPC server.
func Register(gRPC *grpc.Server, authService service.AuthService) {
	ssopb.RegisterSsoServer(gRPC, &serverAPI{auth: authService})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssopb.LoginRequest,
) (*ssopb.LoginResponse, error) {
	// TODO: validate request

	// TODO: call service login method

	// TODO: returns tokens in response

	return &ssopb.LoginResponse{AccessToken: "", RefreshToken: ""}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssopb.RegisterRequest,
) (*ssopb.RegisterResponse, error) {
	// TODO: validate request

	// TODO: call service register method

	// TODO: returns tokens in response

	return &ssopb.RegisterResponse{AccessToken: "", RefreshToken: ""}, nil
}

func (s *serverAPI) RefreshTokens(
	ctx context.Context,
	req *ssopb.RefreshTokensRequest,
) (*ssopb.RefreshTokensResponse, error) {
	// TODO: validate request

	// TODO: call service refresh tokens method

	// TODO: returns tokens in response

	return &ssopb.RefreshTokensResponse{AccessToken: "", RefreshToken: ""}, nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *ssopb.LogoutRequest,
) (*ssopb.LogoutResponse, error) {
	// TODO: validate request

	// TODO: call service logout method

	return &ssopb.LogoutResponse{Success: true}, nil
}
