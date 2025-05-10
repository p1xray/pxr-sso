package validator

import "context"

type CustomClaims interface {
	Validate(ctx context.Context) error
}
