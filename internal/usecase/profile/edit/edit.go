package edit

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

// Repository is a repository for edit user profile data use-case.
type Repository interface {
	UserProfile(ctx context.Context, id int64) (dto.UserProfile, error)
	Save(ctx context.Context, user *entity.User) error
}

// UseCase is a use-case for editing user profile data.
type UseCase struct {
	log  *slog.Logger
	repo Repository
}

// New returns new edit user profile data use-case.
func New(log *slog.Logger, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

// Execute executes the use-case for editing user profile data.
func (uc *UseCase) Execute(ctx context.Context, data Params) error {
	const op = "usecase.profile.edit"

	log := uc.log.With(
		slog.String("op", op),
		slog.Int64("user ID", data.ID),
	)

	storageUserData, err := uc.repo.UserProfile(ctx, data.ID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, usecase.ErrUserNotFound)
		}

		log.Error("error getting user profile data", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	user := entity.NewUser(
		storageUserData.Username,
		storageUserData.FullName,
		storageUserData.DateOfBirth,
		storageUserData.Gender,
		storageUserData.AvatarFileKey,
		entity.WithUserID(storageUserData.ID),
		entity.WithUserPasswordHash(storageUserData.PasswordHash))

	userUpdateProfileParams := entity.UserUpdateProfileParams{
		FullName:      data.FullName,
		DateOfBirth:   data.DateOfBirth,
		Gender:        data.Gender,
		AvatarFileKey: data.AvatarFileKey,
	}
	user.UpdateProfile(userUpdateProfileParams)

	if err = uc.repo.Save(ctx, &user); err != nil {
		log.Error("error saving user", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
