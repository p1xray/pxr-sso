package oauth

import (
	"context"
	oauthpb "github.com/p1xray/pxr-sso-protos/gen/go/oauth"
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/authorize"
	"google.golang.org/grpc"
)

type serverAPI struct {
	oauthpb.UnimplementedOauthServer
	authorizeUseCase controller.Authorize
}

// RegisterOAuthServer registers the implementation of the API service with the gRPC server.
func RegisterOAuthServer(
	server *grpc.Server,
	authorizeUseCase controller.Authorize,
) {
	api := &serverAPI{
		authorizeUseCase: authorizeUseCase,
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
