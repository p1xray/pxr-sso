package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/internal/lib/logger/sl"
	"log/slog"
)

type Profile struct {
	log     *slog.Logger
	storage ProfileStorage
}

type ProfileStorage interface {
	User(ctx context.Context, id int64) (models.User, error)
}

func NewProfileRepository(log *slog.Logger, storage ProfileStorage) *Profile {
	return &Profile{
		log:     log,
		storage: storage,
	}
}

func (p *Profile) UserProfile(ctx context.Context, id int64) (dto.UserProfile, error) {
	const op = "repository.profile.UserProfile"

	log := p.log.With(
		slog.String("op", op),
		slog.Int64("user ID", id),
	)

	user, err := p.storage.User(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found in storage", sl.Err(err))
		} else {
			log.Error("error getting user profile data from storage", sl.Err(err))
		}

		return dto.UserProfile{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := converter.ToUserProfileDTO(user)

	return userDTO, nil
}
