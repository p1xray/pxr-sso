package storage

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
)

type Storage interface {
	UserByUsername(ctx context.Context, username string) (models.User, error)

	PermissionsByUserID(ctx context.Context, userID int64) ([]models.Permission, error)

	UserClientByCode(ctx context.Context, userID int64, clientCode string) (models.Client, error)
}
