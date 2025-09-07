package card

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/usecase"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

// Repository is a repository for user profile card use-case.
type Repository interface {
	UserProfile(ctx context.Context, id int64) (dto.UserProfile, error)
}

// UseCase is a use-case for getting user profile data.
type UseCase struct {
	log  *slog.Logger
	repo Repository
}

// New returns new user profile card use-case.
func New(log *slog.Logger, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

// Execute executes the use-case for getting user profile data.
func (uc *UseCase) Execute(ctx context.Context, id int64) (entity.User, error) {
	const op = "usecase.profile.card"

	log := uc.log.With(
		slog.String("op", op),
		slog.Int64("user ID", id),
	)

	storageUserData, err := uc.repo.UserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return entity.User{}, fmt.Errorf("%s: %w", op, usecase.ErrUserNotFound)
		}

		log.Error("error getting user profile data", sl.Err(err))

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user := entity.NewUser(
		storageUserData.Username,
		storageUserData.FullName,
		storageUserData.DateOfBirth,
		storageUserData.Gender,
		storageUserData.AvatarFileKey,
		entity.WithUserID(storageUserData.ID))

	return user, nil
}
