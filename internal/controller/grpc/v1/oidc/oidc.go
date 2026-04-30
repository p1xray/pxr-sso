package oidc

import (
	"context"
	oauthpb "github.com/p1xray/pxr-sso-protos/gen/go/oauth"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/consent"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/login"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/register"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/token"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"google.golang.org/grpc"
)

// AuthorizeUseCase is the use case which handles the processing of authorization requests by validating and
// then processing these requests based on defined business logic. It also includes
// the fetching of authorization requests when necessary.
type AuthorizeUseCase interface {
	Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, loginRequest dto.LoginRequest) (dto.LoginResponse, error)
}

type RegisterUseCase interface {
	Execute(ctx context.Context, registerRequest dto.RegisterRequest) (dto.RegisterResponse, error)
}

type ConsentUseCase interface {
	Execute(ctx context.Context, consentRequest dto.ConsentRequest) (string, error)
}

type TokenUseCase interface {
	Execute(ctx context.Context, tokenRequest dto.TokenRequest) (dto.TokenResponse, error)
}

// server handles authentication-related processes in the context of OpenID Connect and OAuth2 protocols.
type server struct {
	oauthpb.UnimplementedOauthServer
	authorize       AuthorizeUseCase
	loginUseCase    LoginUseCase
	registerUseCase RegisterUseCase
	consentUseCase  ConsentUseCase
	tokenUseCase    TokenUseCase
}

// RegisterOIDCServer registers the implementation of the OIDC API handlers with the gRPC server.
func RegisterOIDCServer(
	registrar grpc.ServiceRegistrar,
	authorize AuthorizeUseCase,
	loginUseCase LoginUseCase,
	registerUseCase RegisterUseCase,
	consentUseCase ConsentUseCase,
	tokenUseCase TokenUseCase,
) {
	srv := &server{
		authorize:       authorize,
		loginUseCase:    loginUseCase,
		registerUseCase: registerUseCase,
		consentUseCase:  consentUseCase,
		tokenUseCase:    tokenUseCase,
	}

	oauthpb.RegisterOauthServer(registrar, srv)
}

// Authorize handles requests to the authorization endpoint, performing user authentication and
// getting consent for requested scopes.
//
// This endpoint is a key component of the OpenID Connect flow, initiating user authentication and
// consent for access to their information.
func (s *server) Authorize(
	ctx context.Context,
	req *oauthpb.AuthorizeRequest,
) (*oauthpb.AuthorizeResponse, error) {
	authorizeRequest := dto.NewAuthorizeRequest(
		req.GetResponseType(),
		[]string{"login"}, // TODO: get this from proto
		req.GetClientId(),
		req.GetRedirectUri(),
		req.GetCodeChallenge(),
		req.GetCodeChallengeMethod(),
		req.GetState(),
		req.GetAudience(),
		req.GetScope(),
		[]string{"test_session_id"}, // TODO: get this from proto
	)

	redirectURI, err := s.authorize.Execute(ctx, authorizeRequest)
	if err != nil {
		return nil, err
	}

	return &oauthpb.AuthorizeResponse{RedirectUri: redirectURI}, nil
}

// Login is a gRPC handler for OAuth login.
func (s *server) Login(
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
func (s *server) Register(
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
func (s *server) Consent(
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

// Token processes token issuance requests by evaluating authorization grants and
// issuing access, ID and refresh tokens accordingly.
//
// This endpoint is central to the OAuth 2.0 and OpenID Connect framework,
// facilitating the secure issuance of tokens to authenticated clients.
func (s *server) Token(
	ctx context.Context,
	req *oauthpb.TokenRequest,
) (*oauthpb.TokenResponse, error) {
	tokenParams := token.Params{
		GrantType:         req.GetGrantType(),
		ClientID:          req.GetClientId(),
		AuthorizationCode: req.GetCode(),
		RedirectURI:       req.GetRedirectUri(),
		CodeVerifier:      req.GetCodeVerifier(),
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
