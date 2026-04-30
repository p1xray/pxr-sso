package v1

import (
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/oidc"
	"google.golang.org/grpc"
)

// NewRoutes creates a new routes for the gRPC server controller of version 1.
func NewRoutes(
	registrar grpc.ServiceRegistrar,
	authorizeUseCase oidc.AuthorizeUseCase,
	loginUseCase oidc.LoginUseCase,
	registerUseCase oidc.RegisterUseCase,
	consentUseCase oidc.ConsentUseCase,
	tokenUseCase oidc.TokenUseCase,
) {
	oidc.RegisterOIDCServer(
		registrar,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase)
}
