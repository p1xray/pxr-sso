package grpc

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	v1 "github.com/p1xray/pxr-sso/internal/controller/grpc/v1"
	"google.golang.org/grpc"
)

// NewRouter creates a new router for the gRPC server controller.
func NewRouter(
	server *grpc.Server,
	oldLoginUseCase controller.OldLogin,
	oldRegisterUseCase controller.OldRegister,
	oldRefreshUseCase controller.OldRefreshTokens,
	oldLogoutUseCase controller.OldLogout,
	profileUseCase controller.UserProfile,
	editProfileUseCase controller.EditProfile,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	consentUseCase controller.Consent,
	tokenUseCase controller.Token,
) {
	v1.NewRoutes(
		server,
		oldLoginUseCase,
		oldRegisterUseCase,
		oldRefreshUseCase,
		oldLogoutUseCase,
		profileUseCase,
		editProfileUseCase,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase,
	)
}
