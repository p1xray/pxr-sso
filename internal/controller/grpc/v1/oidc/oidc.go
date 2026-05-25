package oidc

import (
	"context"
	oidcpb "github.com/p1xray/pxr-sso-protos/gen/go/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"google.golang.org/grpc"
)

// Authorize is the handler for authorization requests, ensuring they are processed according to OAuth 2.0
// and OpenID Connect protocol specifications.
type Authorize interface {
	// Execute processes an authorization request.
	Execute(ctx context.Context, req dto.AuthorizeRequest) (string, error)
}

// Token
type Token interface {
	// Execute
	Execute(ctx context.Context, tokenRequest dto.TokenRequest) (dto.TokenResponse, error)
}

// server handles authentication-related processes in the context of OpenID Connect and OAuth2 protocols.
type server struct {
	oidcpb.UnimplementedOidcServer
	authorize Authorize
	token     Token
}

// RegisterOIDCServer registers the implementation of the OIDC API handlers with the gRPC server.
func RegisterOIDCServer(
	registrar grpc.ServiceRegistrar,
	authorize Authorize,
	token Token,
) {
	srv := &server{
		authorize: authorize,
		token:     token,
	}

	oidcpb.RegisterOidcServer(registrar, srv)
}

// Authorize handles requests to the authorization endpoint, performing user authentication and
// getting consent for requested scopes.
//
// This endpoint is a key component of the OpenID Connect flow, initiating user authentication and
// consent for access to their information.
func (s *server) Authorize(
	ctx context.Context,
	req *oidcpb.AuthorizeRequest,
) (*oidcpb.AuthorizeResponse, error) {
	sessionCookies := make([]dto.SessionCookie, len(req.GetSessions()))
	for i, session := range req.GetSessions() {
		sessionCookies[i] = dto.NewSessionCookie(session.GetName(), session.GetValue())
	}

	authorizeRequest := dto.NewAuthorizeRequest(
		req.GetResponseType(),
		req.GetPrompt(),
		req.GetClientId(),
		req.GetRedirectUri(),
		req.GetCodeChallenge(),
		req.GetCodeChallengeMethod(),
		req.GetState(),
		req.GetAudience(),
		req.GetScope(),
		sessionCookies,
	)

	authorizeResponse, err := s.authorize.Execute(ctx, authorizeRequest)
	if err != nil {
		return nil, err
	}

	return &oidcpb.AuthorizeResponse{RedirectUri: authorizeResponse}, nil
}

// Token is a gRPC handler for OAuth exchange token.

// Token processes token issuance requests by evaluating authorization grants and
// issuing access, ID and refresh tokens accordingly.
//
// This endpoint is central to the OAuth 2.0 and OpenID Connect framework,
// facilitating the secure issuance of tokens to authenticated clients.
func (s *server) Token(
	ctx context.Context,
	req *oidcpb.TokenRequest,
) (*oidcpb.TokenResponse, error) {
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

	return &oidcpb.TokenResponse{
		AccessToken:  tokenResponse.AccessToken(),
		TokenType:    tokenResponse.TokenType(),
		ExpiresIn:    tokenResponse.ExpiresIn(),
		RefreshToken: tokenResponse.RefreshToken(),
		IdToken:      tokenResponse.IDToken(),
	}, nil
}
