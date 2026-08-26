package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/entity"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/enum"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/models"
)

type SessionOption func(context.Context, *repository, *models.Session) error

func (r *repository) Session(ctx context.Context, id int64, opts ...SessionOption) (dto.Session, error) {
	const op = "get session by id"

	session, err := r.storage.Session(ctx, id)
	if err != nil {
		return dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	for _, opt := range opts {
		if err = opt(ctx, r, &session); err != nil {
			return dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}
	}

	sessionDTO := converter.ToSessionDTO(session)
	return sessionDTO, nil
}

func (r *repository) SessionsByCode(ctx context.Context, codes []string, opts ...SessionOption) ([]dto.Session, error) {
	const op = "get sessions by code"

	sessions, err := r.storage.SessionsByCode(ctx, codes)
	if err != nil {
		return []dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	for _, opt := range opts {
		for i := range sessions {
			if err = opt(ctx, r, &sessions[i]); err != nil {
				return []dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
			}
		}
	}

	sessionsDTO := converter.ToSessionsDTO(sessions)
	return sessionsDTO, nil
}

func (r *repository) SessionByCode(ctx context.Context, code string, opts ...SessionOption) (dto.Session, error) {
	const op = "get session by code"

	session, err := r.storage.SessionByCode(ctx, code)
	if err != nil {
		return dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	for _, opt := range opts {
		if err = opt(ctx, r, &session); err != nil {
			return dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
		}
	}

	sessionDTO := converter.ToSessionDTO(session)
	return sessionDTO, nil
}

func (r *repository) sessionUser(ctx context.Context, userID int64) (models.User, error) {
	user, err := r.storage.User(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *repository) sessionClient(ctx context.Context, clientID int64) (models.Client, error) {
	client, err := r.storage.Client(ctx, clientID)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}

func (r *repository) sessionGrantedScopes(ctx context.Context, sessionID int64) ([]models.SessionGrantedScopeLink, error) {
	links, err := r.storage.SessionGrantedScopeLinks(ctx, sessionID)
	if err != nil {
		return []models.SessionGrantedScopeLink{}, err
	}

	scopeIDs := make([]int64, len(links))
	for i, link := range links {
		scopeIDs[i] = link.ScopeID
	}

	scopes, err := r.storage.Scopes(ctx, scopeIDs)
	if err != nil {
		return []models.SessionGrantedScopeLink{}, err
	}

	for i := range links {
		for j := range scopes {
			if links[i].ScopeID == scopes[j].ID {
				links[i].Scope = scopes[j]
				break
			}
		}
	}

	return links, nil
}

func (r *repository) SaveSession(ctx context.Context, session entity.Session) (dto.Session, error) {
	const op = "save session"

	savedSession, err := r.saveSessionByDataStatus(ctx, session)
	if err != nil {
		return dto.Session{}, fmt.Errorf("%s: %s: %w", pkgTag, op, err)
	}

	return savedSession, nil
}

func (r *repository) saveSessionByDataStatus(ctx context.Context, session entity.Session) (dto.Session, error) {
	switch session.DataStatus() {
	case enum.DataStatusToCreate:
		sessionID, err := r.createSession(ctx, session)
		if err != nil {
			return dto.Session{}, err
		}

		savedSession, err := r.Session(ctx, sessionID, WithUser(), WithClient(), WithGrantedScopes())
		if err != nil {
			return dto.Session{}, err
		}

		return savedSession, nil

	case enum.DataStatusToUpdate:
		if err := r.updateSession(ctx, session); err != nil {
			return dto.Session{}, err
		}

		savedSession, err := r.Session(ctx, session.ID(), WithUser(), WithClient(), WithGrantedScopes())
		if err != nil {
			return dto.Session{}, err
		}

		return savedSession, nil
	case enum.DataStatusToRemove:
		if err := r.removeSession(ctx, session); err != nil {
			return dto.Session{}, err
		}

		return dto.Session{}, nil
	default:
		return dto.Session{}, fmt.Errorf(
			"there is no implementation of save session for data status with value: %d", session.DataStatus())
	}
}

func (r *repository) createSession(ctx context.Context, session entity.Session) (int64, error) {
	sessionStorageModel := converter.ToSessionStorage(session, models.SessionCreated())

	newSessionID := int64(0)
	err := r.storage.WithTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		newSessionID, err = r.storage.CreateSession(ctx, tx, sessionStorageModel)
		if err != nil {
			return err
		}

		err = r.actualizeSessionGrantedScopeLinks(
			ctx,
			tx,
			[]models.SessionGrantedScopeLink{},
			newSessionID,
			session.Scopes(),
		)
		if err != nil {
			return err
		}

		return nil
	})

	return newSessionID, err
}

func (r *repository) updateSession(ctx context.Context, session entity.Session) error {
	sessionStorageModel := converter.ToSessionStorage(session, models.SessionUpdated())

	currentSessionGrantedScopeLinks, err := r.storage.SessionGrantedScopeLinks(ctx, session.ID())
	if err != nil {
		return err
	}

	err = r.storage.WithTransaction(ctx, func(tx pgx.Tx) error {
		if err := r.storage.UpdateSession(ctx, tx, sessionStorageModel); err != nil {
			return err
		}

		err = r.actualizeSessionGrantedScopeLinks(
			ctx,
			tx,
			currentSessionGrantedScopeLinks,
			session.ID(),
			session.Scopes(),
		)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *repository) removeSession(ctx context.Context, session entity.Session) error {
	err := r.storage.WithTransaction(ctx, func(tx pgx.Tx) error {
		if err := r.storage.RemoveSessionGrantedScopeLinksBySessionID(ctx, tx, session.ID()); err != nil {
			return err
		}

		if err := r.storage.RemoveSession(ctx, tx, session.ID()); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *repository) actualizeSessionGrantedScopeLinks(
	ctx context.Context,
	tx pgx.Tx,
	currentLinks []models.SessionGrantedScopeLink,
	sessionID int64,
	scopes []dto.Scope,
) error {
	linksToRemove := make([]int64, 0)
	for _, currentLink := range currentLinks {
		found := false
		for _, scope := range scopes {
			if scope.ID() == currentLink.ScopeID {
				found = true
			}
		}

		if !found {
			linksToRemove = append(linksToRemove, currentLink.ID)
		}
	}

	if err := r.storage.RemoveSessionGrantedScopeLinks(ctx, tx, linksToRemove); err != nil {
		return err
	}

	linksToCreate := make([]models.SessionGrantedScopeLink, 0)
	for _, scope := range scopes {
		found := false
		for _, currentLink := range currentLinks {
			if currentLink.ScopeID == scope.ID() {
				found = true
			}
		}

		if !found {
			linkToCreate := converter.ToSessionGrantedScopeLinkStorage(
				sessionID,
				scope.ID(),
				models.SessionGrantedScopeLinkCreated())
			linksToCreate = append(linksToCreate, linkToCreate)
		}
	}

	if err := r.storage.CreateSessionGrantedScopeLinks(ctx, tx, linksToCreate); err != nil {
		return err
	}

	return nil
}

func WithUser() SessionOption {
	return func(ctx context.Context, r *repository, session *models.Session) error {
		user, err := r.sessionUser(ctx, session.UserID)
		if err != nil {
			return err
		}

		session.User = user
		return nil
	}
}

func WithClient() SessionOption {
	return func(ctx context.Context, r *repository, session *models.Session) error {
		client, err := r.sessionClient(ctx, session.ClientID)
		if err != nil {
			return err
		}

		session.Client = client
		return nil
	}
}

func WithGrantedScopes() SessionOption {
	return func(ctx context.Context, r *repository, session *models.Session) error {
		links, err := r.sessionGrantedScopes(ctx, session.ID)
		if err != nil {
			return err
		}

		session.GrantedScopeLinks = links
		return nil
	}
}
