package grpc

import (
	v1 "github.com/p1xray/pxr-sso/internal/controller/grpc/v1"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/oidc"
	"google.golang.org/grpc"
)

// NewRouter creates a new router for the gRPC server controller.
func NewRouter(
	registrar grpc.ServiceRegistrar,
	authorizeUseCase oidc.Authorize,
	loginUseCase oidc.Login,
	registerUseCase oidc.Register,
	consentUseCase oidc.Consent,
	tokenUseCase oidc.Token,
) {
	v1.NewRoutes(
		registrar,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase,
	)
}
