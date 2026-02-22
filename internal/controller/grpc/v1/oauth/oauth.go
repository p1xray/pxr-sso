package oauth

import (
	"context"
	oauthpb "github.com/p1xray/pxr-sso-protos/gen/go/oauth"
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/consent"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/login"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/register"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/token"
	"google.golang.org/grpc"
)

type serverAPI struct {
	oauthpb.UnimplementedOauthServer
	authorizeUseCase controller.Authorize
	loginUseCase     controller.Login
	registerUseCase  controller.Register
	consentUseCase   controller.Consent
	tokenUseCase     controller.Token
}

// RegisterOAuthServer registers the implementation of the API service with the gRPC server.
func RegisterOAuthServer(
	server *grpc.Server,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	consentUseCase controller.Consent,
	tokenUseCase controller.Token,
) {
	api := &serverAPI{
		authorizeUseCase: authorizeUseCase,
		loginUseCase:     loginUseCase,
		registerUseCase:  registerUseCase,
		consentUseCase:   consentUseCase,
		tokenUseCase:     tokenUseCase,
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
		Scope:               req.GetScope(),
	}

	authorizeData, err := s.authorizeUseCase.Execute(ctx, authorizeParams)
	response := &oauthpb.AuthorizeResponse{
		RedirectUri: authorizeData.Data(),
		
		Error: &oauthpb.ErrorResponse{
			Code:        authorizeData.ErrorCode(),
			Description: authorizeData.ErrorDescription(),
			Uri:         authorizeData.ErrorURI(),
		},
	}

	return response, err
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

	loginData, err := s.loginUseCase.Execute(ctx, loginParams)
	response := &oauthpb.LoginResponse{
		RedirectUri: loginData.Data(),

		Error: &oauthpb.ErrorResponse{
			Code:        loginData.ErrorCode(),
			Description: loginData.ErrorDescription(),
			Uri:         loginData.ErrorURI(),
		},
	}

	return response, err
}

// Register is a gRPC handler for OAuth register.
func (s *serverAPI) Register(
	ctx context.Context,
	req *oauthpb.RegisterRequest,
) (*oauthpb.RegisterResponse, error) {
	registerParams := register.Params{
		FlowID:       req.GetFlowId(),
		ResponseType: req.GetResponseType(),
		ClientID:     req.GetClientId(),
		RedirectURI:  req.GetRedirectUri(),
		State:        req.GetState(),
		Scope:        req.GetScope(),
		Username:     req.GetUsername(),
		Password:     req.GetPassword(),
		FullName:     req.GetFullName(),
	}

	registerData, err := s.registerUseCase.Execute(ctx, registerParams)
	response := &oauthpb.RegisterResponse{
		RedirectUri: registerData.Data(),

		Error: &oauthpb.ErrorResponse{
			Code:        registerData.ErrorCode(),
			Description: registerData.ErrorDescription(),
			Uri:         registerData.ErrorURI(),
		},
	}

	return response, err
}

// Consent is a gRPC handler for OAuth confirming consent.
func (s *serverAPI) Consent(
	ctx context.Context,
	req *oauthpb.ConsentRequest,
) (*oauthpb.ConsentResponse, error) {
	consentParams := consent.Params{
		FlowID:       req.GetFlowId(),
		ResponseType: req.GetResponseType(),
		ClientID:     req.GetClientId(),
		RedirectURI:  req.GetRedirectUri(),
		State:        req.GetState(),
		Scope:        req.GetScope(),
	}

	consentData, err := s.consentUseCase.Execute(ctx, consentParams)
	response := &oauthpb.ConsentResponse{
		RedirectUri: consentData.Data(),

		Error: &oauthpb.ErrorResponse{
			Code:        consentData.ErrorCode(),
			Description: consentData.ErrorDescription(),
			Uri:         consentData.ErrorURI(),
		},
	}

	return response, err
}

// Token is a gRPC handler for OAuth exchange token.
func (s *serverAPI) Token(
	ctx context.Context,
	req *oauthpb.TokenRequest,
) (*oauthpb.TokenResponse, error) {
	tokenParams := token.Params{
		FlowID:            req.GetFlowId(),
		GrantType:         req.GetGrantType(),
		ClientID:          req.GetClientId(),
		AuthorizationCode: req.GetCode(),
		RedirectURI:       req.GetRedirectUri(),
		CodeVerifier:      req.GetCodeVerifier(),
		// TODO: add audience to proto
		// Audience:          req.GetAudience(),
		Scope: req.GetScope(),
	}

	tokenData, err := s.tokenUseCase.Execute(ctx, tokenParams)

	tokens := tokenData.Data()
	response := &oauthpb.TokenResponse{
		AccessToken:  tokens.AccessToken(),
		TokenType:    tokens.TokenType(),
		ExpiresIn:    tokens.ExpiresIn(),
		RefreshToken: tokens.RefreshToken(),
		IdToken:      tokens.IDToken(),

		Error: &oauthpb.ErrorResponse{
			Code:        tokenData.ErrorCode(),
			Description: tokenData.ErrorDescription(),
			Uri:         tokenData.ErrorURI(),
		},
	}

	return response, err
}
