package jwtparser

import (
	"fmt"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	jwtmiddleware "github.com/p1xray/pxr-sso/pkg/jwt"
)

// ParseAccessToken parses access token using a secret key into a set of claims.
func ParseAccessToken(
	token *jwt.JSONWebToken,
	secretKey []byte,
	customClaimsFunc func() jwtmiddleware.CustomClaims,
) (jwtmiddleware.AccessTokenClaims, jwtmiddleware.CustomClaims, error) {
	claims := []interface{}{&jwt.Claims{}, &jwtmiddleware.RegisteredCustomClaims{}}
	if customClaimsExist(customClaimsFunc) {
		claims = append(claims, customClaimsFunc())
	}

	if err := token.Claims(secretKey, &claims); err != nil {
		return jwtmiddleware.AccessTokenClaims{}, nil, fmt.Errorf("error getting token claims: %w", err)
	}

	defaultClaims := *claims[0].(*jwt.Claims)
	registeredCustomClaims := *claims[1].(*jwtmiddleware.RegisteredCustomClaims)

	registeredClaims := jwtmiddleware.AccessTokenClaims{
		Claims:                 defaultClaims,
		RegisteredCustomClaims: registeredCustomClaims,
	}

	var customClaims jwtmiddleware.CustomClaims
	if len(claims) > 2 {
		customClaims = claims[2].(jwtmiddleware.CustomClaims)
	}

	return registeredClaims, customClaims, nil
}

// ParseRefreshToken parses refresh token as a string using a secret key into a set of claims.
func ParseRefreshToken(tokenStr string, secretKey []byte) (jwtmiddleware.RefreshTokenClaims, error) {
	token, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return jwtmiddleware.RefreshTokenClaims{}, err
	}

	claims := jwtmiddleware.RefreshTokenClaims{}
	if err = token.Claims(secretKey, &claims); err != nil {
		return jwtmiddleware.RefreshTokenClaims{}, err
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

func customClaimsExist(customClaimsFunc func() jwtmiddleware.CustomClaims) bool {
	return customClaimsFunc != nil && customClaimsFunc() != nil
}
