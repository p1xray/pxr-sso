package v1

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/auth"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/profile"
	"google.golang.org/grpc"
)

// NewRoutes creates a new routes for the gRPC server controller of version 1.
func NewRoutes(
	server *grpc.Server,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	refreshUseCase controller.RefreshTokens,
	logoutUseCase controller.Logout,
	profileUseCase controller.UserProfile,
) {
	auth.RegisterAuthServer(
		server,
		loginUseCase,
		registerUseCase,
		refreshUseCase,
		logoutUseCase)

	profile.RegisterProfileServer(server, profileUseCase)
}
