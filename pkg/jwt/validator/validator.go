package validator

import (
	"context"
	"errors"
	
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

type Validator struct {
	keyFunc        func(context.Context) ([]byte, error)
	expectedClaims jwt.Expected
	customClaims   func() CustomClaims
}

func New(
	keyFunc func(context.Context) ([]byte, error),
	issuer string,
	audience []string,
	options ...Option,
) (*Validator, error) {
	if keyFunc == nil {
		return nil, errors.New("keyFunc is required")
	}

	if issuer == "" {
		return nil, errors.New("issuer is required")
	}

	if len(audience) == 0 {
		return nil, errors.New("audience is required")
	}

	validator := &Validator{
		keyFunc: keyFunc,
		expectedClaims: jwt.Expected{
			Issuer:   issuer,
			Audience: audience,
		},
	}

	for _, opt := range options {
		opt(validator)
	}

	return validator, nil
}
