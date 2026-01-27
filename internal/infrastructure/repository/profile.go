package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/dto"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/models"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

type Profile struct {
	log     *slog.Logger
	storage ProfileStorage
}

type ProfileStorage interface {
	User(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
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

func (p *Profile) Save(ctx context.Context, user *entity.User) error {
	const op = "repository.profile.Save"

	log := p.log.With(
		slog.String("op", op),
	)

	if user.IsToCreate() || user.IsToRemove() {
		log.Warn("Attempt to create or delete a user from the profile card")

		return fmt.Errorf("%s: %w", op, infrastructure.ErrCreateOrRemoveUserFromProfileCard)
	}

	if user.IsToUpdate() {
		if err := p.updateUser(ctx, user); err != nil {
			log.Error("error updating user", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (p *Profile) updateUser(ctx context.Context, user *entity.User) error {
	if user.ID == emptyID {
		return infrastructure.ErrRequireIDToUpdate
	}

	userStorageModel, err := p.storage.User(ctx, user.ID)
	if err != nil {
		return err
	}

	userStorageModel = converter.ToUserStorage(userStorageModel, user, models.UserUpdated())

	if err = p.storage.UpdateUser(ctx, userStorageModel); err != nil {
		return err
	}

	user.ResetDataStatus()

	return nil
}
