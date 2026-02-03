package v1

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/auth"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/oauth"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/v1/profile"
	"google.golang.org/grpc"
)

// NewRoutes creates a new routes for the gRPC server controller of version 1.
func NewRoutes(
	server *grpc.Server,
	oldLoginUseCase controller.OldLogin,
	oldRegisterUseCase controller.OldRegister,
	oldRefreshUseCase controller.OldRefreshTokens,
	oldLogoutUseCase controller.OldLogout,
	profileUseCase controller.UserProfile,
	editProfileUseCase controller.EditProfile,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
) {
	auth.RegisterAuthServer(
		server,
		oldLoginUseCase,
		oldRegisterUseCase,
		oldRefreshUseCase,
		oldLogoutUseCase,
	)

	profile.RegisterProfileServer(server, profileUseCase, editProfileUseCase)

	oauth.RegisterOAuthServer(server, authorizeUseCase, loginUseCase)
}
