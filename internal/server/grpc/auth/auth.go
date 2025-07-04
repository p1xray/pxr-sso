package authserver

import (
	"context"
	"errors"
	ssopb "github.com/p1xray/pxr-sso-protos/gen/go/sso"
	"github.com/p1xray/pxr-sso/internal/logic/dto"
	"github.com/p1xray/pxr-sso/internal/logic/service"
	"github.com/p1xray/pxr-sso/internal/server"
	"google.golang.org/grpc"
	"time"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	ssopb.UnimplementedSsoServer
	auth server.AuthService
}

// Register registers the implementation of the API service with the gRPC server.
func Register(gRPC *grpc.Server, authService server.AuthService) {
	ssopb.RegisterSsoServer(gRPC, &serverAPI{auth: authService})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssopb.LoginRequest,
) (*ssopb.LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	loginData := dto.LoginDTO{
		Username:    req.GetUsername(),
		Password:    req.GetPassword(),
		ClientCode:  req.GetClientCode(),
		UserAgent:   req.GetUserAgent(),
		Fingerprint: req.GetFingerprint(),
		Issuer:      req.GetIssuer(),
	}

	tokens, err := s.auth.Login(ctx, loginData)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, server.InvalidArgumentError("invalid username or password")
		}

		return nil, server.InternalError("failed to login")
	}

	return &ssopb.LoginResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func validateLoginRequest(req *ssopb.LoginRequest) error {
	if req.GetUsername() == "" {
		return server.InvalidArgumentError("username is empty")
	}

	if req.GetPassword() == "" {
		return server.InvalidArgumentError("password is empty")
	}

	if req.GetClientCode() == "" {
		return server.InvalidArgumentError("client code is empty")
	}

	if req.GetUserAgent() == "" {
		return server.InvalidArgumentError("user agent is empty")
	}

	if req.GetFingerprint() == "" {
		return server.InvalidArgumentError("fingerprint is empty")
	}

	if req.GetIssuer() == "" {
		return server.InvalidArgumentError("issuer is empty")
	}

	return nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssopb.RegisterRequest,
) (*ssopb.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	var dateOfBirth *time.Time
	if req.GetDateOfBirth() != nil {
		dateOfBirthPbAsTime := req.GetDateOfBirth().AsTime()
		dateOfBirth = &dateOfBirthPbAsTime
	}

	var gender *dto.GenderEnum
	if req.GetGender() != emptyValue {
		genderEnum := dto.GenderEnum(req.GetGender().Number())
		gender = &genderEnum
	}

	var avatarFileKey *string
	if req.GetAvatarFileKey() != nil {
		avatarFileKeyPbString := req.GetAvatarFileKey().GetValue()
		avatarFileKey = &avatarFileKeyPbString
	}

	registerData := dto.RegisterDTO{
		Username:      req.GetUsername(),
		Password:      req.GetPassword(),
		ClientCode:    req.GetClientCode(),
		FIO:           req.GetFio(),
		DateOfBirth:   dateOfBirth,
		Gender:        gender,
		AvatarFileKey: avatarFileKey,
		UserAgent:     req.GetUserAgent(),
		Fingerprint:   req.GetFingerprint(),
		Issuer:        req.GetIssuer(),
	}

	tokens, err := s.auth.Register(ctx, registerData)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, server.InvalidArgumentError("user with this username already exists")
		}

		return nil, server.InternalError("failed to register")
	}

	return &ssopb.RegisterResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func validateRegisterRequest(req *ssopb.RegisterRequest) error {
	if req.GetUsername() == "" {
		return server.InvalidArgumentError("username is empty")
	}

	if req.GetPassword() == "" {
		return server.InvalidArgumentError("password is empty")
	}

	if req.GetClientCode() == "" {
		return server.InvalidArgumentError("client code is empty")
	}

	if req.GetFio() == "" {
		return server.InvalidArgumentError("FIO is empty")
	}

	if req.GetUserAgent() == "" {
		return server.InvalidArgumentError("user agent is empty")
	}

	if req.GetFingerprint() == "" {
		return server.InvalidArgumentError("fingerprint is empty")
	}

	if req.GetIssuer() == "" {
		return server.InvalidArgumentError("issuer is empty")
	}

	return nil
}

func (s *serverAPI) RefreshTokens(
	ctx context.Context,
	req *ssopb.RefreshTokensRequest,
) (*ssopb.RefreshTokensResponse, error) {
	if err := validateRefreshTokensRequest(req); err != nil {
		return nil, err
	}

	refreshTokensData := dto.RefreshTokensDTO{
		RefreshToken: req.GetRefreshToken(),
		ClientCode:   req.GetClientCode(),
		UserAgent:    req.GetUserAgent(),
		Fingerprint:  req.GetFingerprint(),
		Issuer:       req.GetIssuer(),
	}

	tokens, err := s.auth.RefreshTokens(ctx, refreshTokensData)
	if err != nil {
		return nil, server.InternalError("failed to refresh tokens")
	}

	return &ssopb.RefreshTokensResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func validateRefreshTokensRequest(req *ssopb.RefreshTokensRequest) error {
	if req.GetRefreshToken() == "" {
		return server.InvalidArgumentError("refresh token is empty")
	}

	if req.GetUserAgent() == "" {
		return server.InvalidArgumentError("user agent is empty")
	}

	if req.GetFingerprint() == "" {
		return server.InvalidArgumentError("fingerprint is empty")
	}

	if req.GetClientCode() == "" {
		return server.InvalidArgumentError("client code is empty")
	}

	if req.GetIssuer() == "" {
		return server.InvalidArgumentError("issuer is empty")
	}

	return nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *ssopb.LogoutRequest,
) (*ssopb.LogoutResponse, error) {
	if err := validateLogoutRequest(req); err != nil {
		return &ssopb.LogoutResponse{Success: false}, err
	}

	logoutData := dto.LogoutDTO{
		RefreshToken: req.GetRefreshToken(),
		ClientCode:   req.GetClientCode(),
	}
	if err := s.auth.Logout(ctx, logoutData); err != nil {
		return &ssopb.LogoutResponse{Success: false}, server.InternalError("failed to logout")
	}

	return &ssopb.LogoutResponse{Success: true}, nil
}

func validateLogoutRequest(req *ssopb.LogoutRequest) error {
	if req.GetRefreshToken() == "" {
		return server.InvalidArgumentError("refresh token is empty")
	}

	if req.GetClientCode() == "" {
		return server.InvalidArgumentError("client code is empty")
	}

	return nil
}
