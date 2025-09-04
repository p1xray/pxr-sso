package grpc

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	v1 "github.com/p1xray/pxr-sso/internal/controller/grpc/v1"
	"google.golang.org/grpc"
)

// NewRouter creates a new router for the gRPC server controller.
func NewRouter(
	server *grpc.Server,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	refreshUseCase controller.RefreshTokens,
	logoutUseCase controller.Logout,
	profileUseCase controller.UserProfile,
) {
	v1.NewRoutes(
		server,
		loginUseCase,
		registerUseCase,
		refreshUseCase,
		logoutUseCase,
		profileUseCase)
}
