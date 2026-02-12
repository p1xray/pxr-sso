package v1

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/oauth"
	"google.golang.org/grpc"
)

// NewRoutes creates a new routes for the gRPC server controller of version 1.
func NewRoutes(
	server *grpc.Server,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	consentUseCase controller.Consent,
	tokenUseCase controller.Token,
) {
	oauth.RegisterOAuthServer(
		server,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase)
}
