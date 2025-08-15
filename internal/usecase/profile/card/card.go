package card

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"github.com/p1xray/pxr-sso/internal/usecase"
	"log/slog"
)

type ProfileRepository interface {
	UserProfile(ctx context.Context, id int64) (dto.UserProfile, error)
}

type UseCase struct {
	log  *slog.Logger
	repo ProfileRepository
}

func New(log *slog.Logger, repo ProfileRepository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

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
