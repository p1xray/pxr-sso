package authorize

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

// LoginFlowProcessor is the processor for authorize flow with logging in interaction.
type LoginFlowProcessor interface {
	// AuthorizeWithLogin handles authorize flow with logging in interaction.
	AuthorizeWithLogin(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// NoneFlowProcessor is the processor for authorize without interaction flow.
type NoneFlowProcessor interface {
	// AuthorizeWithoutInteraction handles authorize without interaction flow.
	AuthorizeWithoutInteraction(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// ConsentFlowProcessor is the processor for authorize flow with consent confirming interaction.
type ConsentFlowProcessor interface {
	// AuthorizeWithConsent handles authorize flow with consent confirming interaction.
	AuthorizeWithConsent(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// SelectAccountFlowProcessor is the processor for authorize flow with selecting account interaction.
type SelectAccountFlowProcessor interface {
	// AuthorizeWithSelectAccount handles authorize flow with selecting account interaction.
	AuthorizeWithSelectAccount(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

// switcher is the processor for defining the interaction flow.
type switcher struct {
	login         LoginFlowProcessor
	none          NoneFlowProcessor
	consent       ConsentFlowProcessor
	selectAccount SelectAccountFlowProcessor
}

// NewFlowSwitcher creates a new processor for defining the interaction flow.
func NewFlowSwitcher(
	login LoginFlowProcessor,
	none NoneFlowProcessor,
	consent ConsentFlowProcessor,
	selectAccount SelectAccountFlowProcessor,
) *switcher {
	return &switcher{
		login:         login,
		none:          none,
		consent:       consent,
		selectAccount: selectAccount,
	}
}

// Switch processes the authorize context to decide which interaction flow to execute.
//
// This method analyze authorize request, sessions, granted and pending scopes to decide which interaction flow
// to execute or which error are return.
func (s *switcher) Switch(ctx context.Context, data dto.AuthorizeContext) (string, error) {
	request := data.Request()
	prompt := request.Prompt()
	sessions := data.Sessions()

	switch {
	// if no sessions exist and the prompt forbids user interaction,
	// respond that login is required without allowing user interaction.
	case len(sessions) == 0 && prompt == "none":
		return "", oidc.LoginRequiredError()

	// if multiple sessions exist but the prompt forbids interaction,
	// respond that account selection is required but user interaction is not allowed.
	case len(sessions) > 1 && prompt == "none":
		return "", oidc.AccountSelectionRequiredError()

	// if no sessions exist, or the request explicitly asks for a login, prompt the user for login.
	case len(sessions) == 0 || prompt == "login":
		return s.login.AuthorizeWithLogin(ctx, data)

	// if multiple sessions exist, or the request requires account selection, prompt the user to select an account.
	case len(sessions) > 1 || prompt == "select_account":
		return s.selectAccount.AuthorizeWithSelectAccount(ctx, data)

	// if user interaction is disallowed but consent is necessary, return the consent required error.
	case len(sessions) == 1 && prompt == "none" && len(data.RequestPendingScopes()) > 0:
		return "", oidc.ConsentRequiredError()

	// if consent for required scopes is still pending, handle consent requirements.
	case len(sessions) == 1 && (prompt == "consent" || len(data.RequestPendingScopes()) > 0):
		return s.consent.AuthorizeWithConsent(ctx, data)

	// if consent is not required, handle authorize flow with prompt none.
	case len(sessions) == 1 && prompt == "none" && len(data.RequestPendingScopes()) == 0:
		return s.none.AuthorizeWithoutInteraction(ctx, data)

	// catch any unexpected cases where the session count or prompt state does not match the expected conditions.
	default:
		return "", fmt.Errorf("invalid number of sessions: %d or prompt: %s", len(sessions), prompt)
	}
}
