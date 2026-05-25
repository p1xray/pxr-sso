package oidc

import (
	"context"
	oauthpb "github.com/p1xray/pxr-sso-protos/gen/go/oauth"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"google.golang.org/grpc"
)

// Authorize is the handler for authorization requests, ensuring they are processed according to OAuth 2.0
// and OpenID Connect protocol specifications.
type Authorize interface {
	// Execute processes an authorization request.
	Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error)
}

// Login is the handler for logging in user request.
type Login interface {
	// Execute processes a logging in user request.
	Execute(ctx context.Context, loginRequest dto.LoginRequest) (dto.LoginResponse, error)
}

// Register
type Register interface {
	// Execute
	Execute(ctx context.Context, registerRequest dto.RegisterRequest) (dto.RegisterResponse, error)
}

// Consent
type Consent interface {
	// Execute
	Execute(ctx context.Context, consentRequest dto.ConsentRequest) (string, error)
}

// Token
type Token interface {
	// Execute
	Execute(ctx context.Context, tokenRequest dto.TokenRequest) (dto.TokenResponse, error)
}

// server handles authentication-related processes in the context of OpenID Connect and OAuth2 protocols.
type server struct {
	oauthpb.UnimplementedOauthServer
	authorize Authorize
	login     Login
	register  Register
	consent   Consent
	token     Token
}

// RegisterOIDCServer registers the implementation of the OIDC API handlers with the gRPC server.
func RegisterOIDCServer(
	registrar grpc.ServiceRegistrar,
	authorize Authorize,
	login Login,
	register Register,
	consent Consent,
	token Token,
) {
	srv := &server{
		authorize: authorize,
		login:     login,
		register:  register,
		consent:   consent,
		token:     token,
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
	loginRequest := dto.NewLoginRequest(
		"pxr.sso:par:KJFGHDKJGHFKJDGHFJK", // TODO: get this from proto
		req.GetUsername(),
		req.GetPassword(),
	)

	loginResponse, err := s.login.Execute(ctx, loginRequest)
	if err != nil {
		return nil, err
	}

	response := &oauthpb.LoginResponse{
		RedirectUri: loginResponse.RedirectURI(),
		// Session:     loginResponse.Session(), TODO: add this to proto
	}

	return response, nil
}

// Register is a gRPC handler for OAuth register.
func (s *server) Register(
	ctx context.Context,
	req *oauthpb.RegisterRequest,
) (*oauthpb.RegisterResponse, error) {
	registerRequest := dto.NewRegisterRequest(
		"pxr.sso:par:KJFGHDKJGHFKJDGHFJK", // TODO: get this from proto
		req.GetUsername(),
		req.GetPassword(),
		req.GetFullName(),
	)

	registerResponse, err := s.register.Execute(ctx, registerRequest)
	if err != nil {
		return nil, err
	}

	response := &oauthpb.RegisterResponse{
		RedirectUri: registerResponse.RedirectURI(),
		// Session:     loginResponse.Session(), TODO: add this to proto
	}

	return response, nil
}

// Consent is a gRPC handler for OAuth confirming consent.
func (s *server) Consent(
	ctx context.Context,
	req *oauthpb.ConsentRequest,
) (*oauthpb.ConsentResponse, error) {
	consentRequest := dto.NewConsentRequest(
		"pxr.sso:par:KJFGHDKJGHFKJDGHFJK", // TODO: get this from proto
		req.GetScope(),                    // TODO: rename to scopes in proto
		make([]dto.SessionCookie, 0),      // TODO: get this from proto
	)

	consentResponse, err := s.consent.Execute(ctx, consentRequest)
	if err != nil {
		return nil, err
	}

	response := &oauthpb.ConsentResponse{
		RedirectUri: consentResponse,
	}

	return response, nil
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
	tokenRequest := dto.NewTokenRequest(
		req.GetGrantType(),
		req.GetCode(),
		req.GetRedirectUri(),
		req.GetCodeVerifier(),
		req.GetClientId(),
	)

	tokenResponse, err := s.token.Execute(ctx, tokenRequest)
	if err != nil {
		return nil, err
	}

	response := &oauthpb.TokenResponse{
		AccessToken:  tokenResponse.AccessToken(),
		TokenType:    tokenResponse.TokenType(),
		ExpiresIn:    tokenResponse.ExpiresIn(),
		RefreshToken: tokenResponse.RefreshToken(),
		IdToken:      tokenResponse.IDToken(),
	}

	return response, nil
}
