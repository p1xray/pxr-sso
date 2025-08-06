package usecase

import "context"

type (
	Login interface {
		Execute(ctx context.Context) error
	}
)
