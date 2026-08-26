package repository

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage"
)

const pkgTag = "storage repository"

const emptyID = 0

type Repository interface {
	ClientByCode(ctx context.Context, code string, opts ...ClientOption) (dto.Client, error)

	User(ctx context.Context, id int64, opts ...UserOption) (dto.User, error)
	UserByUsername(ctx context.Context, username string, opts ...UserOption) (dto.User, error)
	IsUserExistByUsername(ctx context.Context, username string) (bool, error)
	SaveUser(ctx context.Context, user entity.User) (dto.User, error)

	SessionsByCode(ctx context.Context, codes []string, opts ...SessionOption) ([]dto.Session, error)
	SessionByCode(ctx context.Context, code string, opts ...SessionOption) (dto.Session, error)
	SaveSession(ctx context.Context, session entity.Session) (dto.Session, error)

	ScopesByCode(ctx context.Context, codes []string) ([]dto.Scope, error)
}

type repository struct {
	storage storage.Storage
}

func New(storage storage.Storage) *repository {
	return &repository{
		storage: storage,
	}
}
