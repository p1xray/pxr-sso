package v1

import (
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/oidc"
	"google.golang.org/grpc"
)

// NewRoutes creates a new routes for the gRPC server controller of version 1.
func NewRoutes(
	registrar grpc.ServiceRegistrar,
	authorizeUseCase oidc.Authorize,
	loginUseCase oidc.Login,
	registerUseCase oidc.Register,
	consentUseCase oidc.Consent,
	tokenUseCase oidc.Token,
) {
	oidc.RegisterOIDCServer(
		registrar,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase)
}
