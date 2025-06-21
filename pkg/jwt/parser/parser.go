package jwtparser

import (
	"fmt"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	jwtclaims "github.com/p1xray/pxr-sso/pkg/jwt/claims"
)

// ParseAccessToken parses access token using a secret key into a set of claims.
func ParseAccessToken(
	token *jwt.JSONWebToken,
	secretKey []byte,
	customClaimsFunc func() jwtclaims.CustomClaims,
) (jwtclaims.AccessTokenClaims, jwtclaims.CustomClaims, error) {
	defaultClaims := jwt.Claims{}
	registeredCustomClaims := jwtclaims.RegisteredCustomClaims{}
	var customClaims jwtclaims.CustomClaims

	if err := token.Claims(secretKey, &defaultClaims, &registeredCustomClaims); err != nil {
		return jwtclaims.AccessTokenClaims{}, nil, fmt.Errorf("error getting token claims: %w", err)
	}

	if customClaimsExist(customClaimsFunc) {
		customClaims = customClaimsFunc()
		if err := token.Claims(secretKey, &customClaims); err != nil {
			return jwtclaims.AccessTokenClaims{}, nil, fmt.Errorf("error getting token custom claims: %w", err)
		}
	}

	registeredClaims := jwtclaims.AccessTokenClaims{
		Claims:                 defaultClaims,
		RegisteredCustomClaims: registeredCustomClaims,
	}

	return registeredClaims, customClaims, nil
}

// ParseRefreshToken parses refresh token as a string using a secret key into a set of claims.
func ParseRefreshToken(tokenStr string, secretKey []byte) (jwtclaims.RefreshTokenClaims, error) {
	token, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return jwtclaims.RefreshTokenClaims{}, err
	}

	claims := jwtclaims.RefreshTokenClaims{}
	if err = token.Claims(secretKey, &claims); err != nil {
		return jwtclaims.RefreshTokenClaims{}, err
	}

	return claims, nil
}

// ParseSignatureAlgorithm parses signature algorithm from token header.
func ParseSignatureAlgorithm(token *jwt.JSONWebToken) (jose.SignatureAlgorithm, error) {
	if len(token.Headers) < 1 {
		return "", fmt.Errorf("token header is empty")
	}

	signatureAlgorithm := token.Headers[0].Algorithm
	if signatureAlgorithm == "" {
		return "", fmt.Errorf("signature algorithm is empty")
	}

	return jose.SignatureAlgorithm(signatureAlgorithm), nil
}

func customClaimsExist(customClaimsFunc func() jwtclaims.CustomClaims) bool {
	return customClaimsFunc != nil && customClaimsFunc() != nil
}
