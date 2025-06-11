package profile

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"pxr-sso/internal/lib/logger/sl"
	usercrud "pxr-sso/internal/logic/crud/user"
	"pxr-sso/internal/logic/dto"
	"pxr-sso/internal/logic/service"
	"pxr-sso/internal/storage"
)

// Profile is service for working with user profile.
type Profile struct {
	log      *slog.Logger
	userCRUD *usercrud.CRUD
}

// New creates a new profile service.
func New(
	log *slog.Logger,
	userCRUD *usercrud.CRUD,
) *Profile {
	return &Profile{
		log:      log,
		userCRUD: userCRUD,
	}
}

// UserProfile returns user profile data.
func (p *Profile) UserProfile(ctx context.Context, userID int64) (dto.UserProfileDTO, error) {
	const op = "profile.UserProfile"

	log := p.log.With(
		slog.String("op", op),
		slog.Int64("user id", userID),
	)
	log.Info("attempting to get user profile")
	
	user, err := p.userCRUD.User(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return dto.UserProfileDTO{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("failed to get user", sl.Err(err))

		return dto.UserProfileDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userProfile := dto.UserProfileDTO{
		UserID:        user.ID,
		Username:      user.Username,
		FIO:           user.FIO,
		DateOfBirth:   user.DateOfBirth,
		Gender:        user.Gender,
		AvatarFileKey: user.AvatarFileKey,
	}

	return userProfile, nil
}
