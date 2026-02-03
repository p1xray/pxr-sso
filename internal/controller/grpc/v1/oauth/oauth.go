package oauth

import (
	"context"
	oauthpb "github.com/p1xray/pxr-sso-protos/gen/go/oauth"
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/login"
	"google.golang.org/grpc"
)

type serverAPI struct {
	oauthpb.UnimplementedOauthServer
	authorizeUseCase controller.Authorize
	loginUseCase     controller.Login
}

// RegisterOAuthServer registers the implementation of the API service with the gRPC server.
func RegisterOAuthServer(
	server *grpc.Server,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
) {
	api := &serverAPI{
		authorizeUseCase: authorizeUseCase,
		loginUseCase:     loginUseCase,
	}

	oauthpb.RegisterOauthServer(server, api)
}

// Authorize is a gRPC handler for OAuth authorize.
func (s *serverAPI) Authorize(
	ctx context.Context,
	req *oauthpb.AuthorizeRequest,
) (*oauthpb.AuthorizeResponse, error) {
	authorizeParams := authorize.Params{
		ResponseType:        req.GetResponseType(),
		ClientID:            req.GetClientId(),
		RedirectURI:         req.GetRedirectUri(),
		CodeChallenge:       req.GetCodeChallenge(),
		CodeChallengeMethod: req.GetCodeChallengeMethod(),
		State:               req.GetState(),
	}
	redirectURI := s.authorizeUseCase.Execute(ctx, authorizeParams)

	response := &oauthpb.AuthorizeResponse{
		RedirectUri: redirectURI,
	}
	return response, nil
}

// Login is a gRPC handler for OAuth login.
func (s *serverAPI) Login(
	ctx context.Context,
	req *oauthpb.LoginRequest,
) (*oauthpb.LoginResponse, error) {
	loginParams := login.Params{
		FlowID:       req.GetFlowId(),
		ResponseType: req.GetResponseType(),
		ClientID:     req.GetClientId(),
		RedirectURI:  req.GetRedirectUri(),
		State:        req.GetState(),
		Scope:        req.GetScope(),
		Username:     req.GetUsername(),
		Password:     req.GetPassword(),
	}
	redirectURI, err := s.loginUseCase.Execute(ctx, loginParams)
	if err != nil {
		// TODO: update proto with displayable error

		return nil, err.Unwrap()
	}

	response := &oauthpb.LoginResponse{
		RedirectUri: redirectURI,
	}
	return response, nil
}
