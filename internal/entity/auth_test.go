package entity

import (
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	accessTokenTTL  = time.Minute
	refreshTokenTTL = time.Minute

	userID              = 1
	validPassword       = "123456"
	invalidPassword     = "1"
	invalidLongPassword = "wJBXKFhPtiCtsWKTRyDWHWMmcPnvwGIjROLxRZyvRNyFktpAYZvidWclgXERJZYHVHjlHwiHY"
	passwordHash        = "$2a$10$MzZyhvwQgTriuJ3pPH0z.exkmUxk2gwV1vFHlQvd9457n27gGR4NO"

	sessionID      = 1
	refreshTokenID = "6424f67d-61c3-4251-b193-f2da172f9e01"

	clientID         = 1
	userAgent        = "test user agent"
	fingerprint      = "test fingerprint"
	issuer           = "test issuer"
	secretKey        = "98649a5c-2137-4a78-a63f-fbab416a7f9e"
	invalidSecretKey = "invalid_key"
)

func Test_Auth_Login(t *testing.T) {
	testCases := []struct {
		name          string
		data          LoginParams
		user          dto.User
		client        dto.Client
		expectedError error
	}{
		{
			name: "successfully log in",
			data: LoginParams{
				Password:    validPassword,
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			user: dto.User{
				ID:           userID,
				PasswordHash: passwordHash,
			},
			client: dto.Client{
				ID:        clientID,
				SecretKey: secretKey,
			},
			expectedError: nil,
		},
		{
			name: "throws an error when user is empty",
			data: LoginParams{
				Password:    validPassword,
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			client: dto.Client{
				ID:        clientID,
				SecretKey: secretKey,
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name: "throws an error when client is empty",
			data: LoginParams{
				Password:    validPassword,
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			user: dto.User{
				ID:           userID,
				PasswordHash: passwordHash,
			},
			expectedError: ErrCreateSession,
		},
		{
			name: "throws an error when given password is invalid",
			data: LoginParams{
				Password:    invalidPassword,
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			user: dto.User{
				ID:           userID,
				PasswordHash: passwordHash,
			},
			client: dto.Client{
				ID:        clientID,
				SecretKey: secretKey,
			},
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth, err := NewAuth(accessTokenTTL, refreshTokenTTL, WithUser(tc.user), WithClient(tc.client))
			require.NoError(t, err)

			tokens, err := auth.Login(tc.data)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.NotEmpty(t, tokens.RefreshTokenID)

				assert.True(t, auth.Sessions[0].IsToCreate())
			}
		})
	}
}

func Test_Auth_Register(t *testing.T) {
	var (
		username      = "test@mail.com"
		fullName      = "test user"
		dateOfBirth   = time.Date(2000, 5, 16, 0, 0, 0, 0, time.Local)
		gender        = enum.MALE
		avatarFileKey = "test avatar file key"
	)

	expectedUser := User{
		Username:      username,
		FullName:      fullName,
		DateOfBirth:   &dateOfBirth,
		Gender:        &gender,
		AvatarFileKey: &avatarFileKey,
	}

	testCases := []struct {
		name          string
		data          RegisterParams
		user          dto.User
		client        dto.Client
		expectedUser  User
		expectedError error
	}{
		{
			name: "successfully register",
			data: RegisterParams{
				Username:      username,
				Password:      validPassword,
				FullName:      fullName,
				DateOfBirth:   &dateOfBirth,
				Gender:        &gender,
				AvatarFileKey: &avatarFileKey,
				Fingerprint:   fingerprint,
				Issuer:        issuer,
			},
			client: dto.Client{
				ID:        clientID,
				SecretKey: secretKey,
			},
			expectedUser:  expectedUser,
			expectedError: nil,
		},
		{
			name: "throws an error when user is already registered",
			data: RegisterParams{
				Username:      username,
				Password:      validPassword,
				FullName:      fullName,
				DateOfBirth:   &dateOfBirth,
				Gender:        &gender,
				AvatarFileKey: &avatarFileKey,
				Fingerprint:   fingerprint,
				Issuer:        issuer,
			},
			user: dto.User{
				ID: userID,
			},
			expectedError: ErrUserExists,
		},
		{
			name: "throws an error when password is too long",
			data: RegisterParams{
				Username:      username,
				Password:      invalidLongPassword,
				FullName:      fullName,
				DateOfBirth:   &dateOfBirth,
				Gender:        &gender,
				AvatarFileKey: &avatarFileKey,
				Fingerprint:   fingerprint,
				Issuer:        issuer,
			},
			expectedError: ErrGeneratePasswordHash,
		},
		{
			name: "throws an error when client is empty",
			data: RegisterParams{
				Username:      username,
				Password:      validPassword,
				FullName:      fullName,
				DateOfBirth:   &dateOfBirth,
				Gender:        &gender,
				AvatarFileKey: &avatarFileKey,
				Fingerprint:   fingerprint,
				Issuer:        issuer,
			},
			expectedError: ErrCreateSession,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth, err := NewAuth(accessTokenTTL, refreshTokenTTL, WithUser(tc.user), WithClient(tc.client))
			require.NoError(t, err)

			tokens, err := auth.Register(tc.data)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.NotEmpty(t, tokens.RefreshTokenID)

				assert.Equal(t, tc.expectedUser.ID, auth.User.ID)
				assert.Equal(t, tc.expectedUser.Username, auth.User.Username)
				assert.Equal(t, tc.expectedUser.FullName, auth.User.FullName)
				assert.Equal(t, tc.expectedUser.DateOfBirth, auth.User.DateOfBirth)
				assert.Equal(t, tc.expectedUser.Gender, auth.User.Gender)
				assert.Equal(t, tc.expectedUser.AvatarFileKey, auth.User.AvatarFileKey)
				assert.True(t, auth.User.IsToCreate())
				assert.True(t, auth.Sessions[0].IsToCreate())
			}
		})
	}
}

func Test_Auth_RefreshTokens(t *testing.T) {
	var (
		validSessionExpires   = time.Now().Add(time.Hour)
		invalidSessionExpires = time.Now().Add(-time.Hour)
	)

	testCases := []struct {
		name          string
		data          RefreshTokensParams
		user          dto.User
		client        dto.Client
		session       dto.Session
		expectedError error
	}{
		{
			name: "successfully refresh tokens",
			data: RefreshTokensParams{
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			user: dto.User{
				ID: userID,
			},
			client: dto.Client{
				ID:        clientID,
				SecretKey: secretKey,
			},
			session: dto.Session{
				ID:             sessionID,
				UserID:         userID,
				RefreshTokenID: refreshTokenID,
				UserAgent:      userAgent,
				Fingerprint:    fingerprint,
				ExpiresAt:      validSessionExpires,
			},
			expectedError: nil,
		},
		{
			name: "throws an error when session is expired",
			data: RefreshTokensParams{
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			session: dto.Session{
				ID:             sessionID,
				UserID:         userID,
				RefreshTokenID: refreshTokenID,
				UserAgent:      userAgent,
				Fingerprint:    fingerprint,
				ExpiresAt:      invalidSessionExpires,
			},
			expectedError: ErrValidateSession,
		},
		{
			name: "throws an error when client is empty",
			data: RefreshTokensParams{
				UserAgent:   userAgent,
				Fingerprint: fingerprint,
				Issuer:      issuer,
			},
			user: dto.User{
				ID: userID,
			},
			session: dto.Session{
				ID:             sessionID,
				UserID:         userID,
				RefreshTokenID: refreshTokenID,
				UserAgent:      userAgent,
				Fingerprint:    fingerprint,
				ExpiresAt:      validSessionExpires,
			},
			expectedError: ErrCreateSession,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth, err := NewAuth(
				accessTokenTTL,
				refreshTokenTTL,
				WithUser(tc.user),
				WithClient(tc.client),
				WithSession(tc.session),
			)
			require.NoError(t, err)

			tokens, err := auth.RefreshTokens(tc.data)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.NotEmpty(t, tokens.RefreshTokenID)

				var previousSession Session
				var newSession Session
				for _, session := range auth.Sessions {
					if session.ID == sessionID {
						previousSession = session
					} else {
						newSession = session
					}
				}

				assert.True(t, previousSession.IsToRemove())
				assert.True(t, newSession.IsToCreate())
			}
		})
	}
}

func Test_Auth_Logout(t *testing.T) {
	sessionExpires := time.Now().Add(time.Hour)

	testCases := []struct {
		name          string
		session       dto.Session
		expectedError error
	}{
		{
			name: "successfully logout",
			session: dto.Session{
				ID:             sessionID,
				UserID:         userID,
				RefreshTokenID: refreshTokenID,
				UserAgent:      userAgent,
				Fingerprint:    fingerprint,
				ExpiresAt:      sessionExpires,
			},
			expectedError: nil,
		},
		{
			name:          "throws an error when session is empty",
			expectedError: ErrSessionNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth, err := NewAuth(accessTokenTTL, refreshTokenTTL, WithSession(tc.session))
			require.NoError(t, err)

			err = auth.Logout()

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				for _, session := range auth.Sessions {
					assert.True(t, session.IsToRemove())
				}
			}
		})
	}
}
