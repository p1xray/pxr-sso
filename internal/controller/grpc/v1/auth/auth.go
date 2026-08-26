package auth

import (
	"context"
	authpb "github.com/p1xray/pxr-sso-protos/gen/go/auth"
	sessionpb "github.com/p1xray/pxr-sso-protos/gen/go/session"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"google.golang.org/grpc"
)

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

// ConsentCardReader
type ConsentCardReader interface {
	// Read
	Read(ctx context.Context, request dto.ConsentCardRequest) (dto.ConsentCardResponse, error)
}

// server handles authentication-related processes in the context of OpenID Connect and OAuth2 protocols.
type server struct {
	authpb.UnimplementedAuthServer
	login             Login
	register          Register
	consent           Consent
	consentCardReader ConsentCardReader
}

// RegisterAuthServer registers the implementation of the auth API handlers with the gRPC server.
func RegisterAuthServer(
	registrar grpc.ServiceRegistrar,
	login Login,
	register Register,
	consent Consent,
	consentCardReader ConsentCardReader,
) {
	srv := &server{
		login:             login,
		register:          register,
		consent:           consent,
		consentCardReader: consentCardReader,
	}

	authpb.RegisterAuthServer(registrar, srv)
}

// Login is a gRPC handler for OAuth login.
func (s *server) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	loginRequest := dto.NewLoginRequest(
		req.GetRequestUri(),
		req.GetUsername(),
		req.GetPassword(),
	)

	loginResponse, err := s.login.Execute(ctx, loginRequest)
	if err != nil {
		return nil, err
	}

	session := loginResponse.Session()
	sessionCookie := &sessionpb.Cookie{
		Name:  session.Name(),
		Value: session.Value(),
	}

	return &authpb.LoginResponse{
		RedirectUri: loginResponse.RedirectURI(),
		Session:     sessionCookie,
	}, nil
}

// Register is a gRPC handler for OAuth register.
func (s *server) Register(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	registerRequest := dto.NewRegisterRequest(
		req.GetRequestUri(),
		req.GetUsername(),
		req.GetPassword(),
		req.GetFullName(),
	)

	registerResponse, err := s.register.Execute(ctx, registerRequest)
	if err != nil {
		return nil, err
	}

	session := registerResponse.Session()
	sessionCookie := &sessionpb.Cookie{
		Name:  session.Name(),
		Value: session.Value(),
	}

	return &authpb.RegisterResponse{
		RedirectUri: registerResponse.RedirectURI(),
		Session:     sessionCookie,
	}, nil
}

// Consent is a gRPC handler for OAuth confirming consent.
func (s *server) Consent(
	ctx context.Context,
	req *authpb.ConsentRequest,
) (*authpb.ConsentResponse, error) {
	sessionCookies := make([]dto.SessionCookie, len(req.GetSessions()))
	for i, session := range req.GetSessions() {
		sessionCookies[i] = dto.NewSessionCookie(session.GetName(), session.GetValue())
	}

	consentRequest := dto.NewConsentRequest(
		req.GetRequestUri(),
		req.GetScopes(),
		sessionCookies,
	)

	consentResponse, err := s.consent.Execute(ctx, consentRequest)
	if err != nil {
		return nil, err
	}

	return &authpb.ConsentResponse{RedirectUri: consentResponse}, nil
}

// GetConsentCard is a gRPC handler for getting consent card data.
func (s *server) GetConsentCard(
	ctx context.Context,
	req *authpb.GetConsentCardRequest,
) (*authpb.GetConsentCardResponse, error) {
	consentCardRequest := dto.NewConsentCardRequest(req.GetRequestUri())
	consentCardResponse, err := s.consentCardReader.Read(ctx, consentCardRequest)
	if err != nil {
		return nil, err
	}

	consentScopes := make([]*authpb.ConsentScope, len(consentCardResponse.Scopes()))
	for i, scope := range consentCardResponse.Scopes() {
		consentScope := &authpb.ConsentScope{
			Code:        scope.Code(),
			Name:        scope.Name(),
			Description: scope.Description(),
			IsGranted:   scope.IsGranted(),
		}

		consentScopes[i] = consentScope
	}

	return &authpb.GetConsentCardResponse{Scopes: consentScopes}, nil
}
