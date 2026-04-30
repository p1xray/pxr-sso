package authorize

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
)

type LoginFlowProcessor interface {
	AuthorizeWithLogin(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type NoneFlowProcessor interface {
	AuthorizeWithoutInteraction(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type ConsentFlowProcessor interface {
	AuthorizeWithConsent(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type SelectAccountFlowProcessor interface {
	AuthorizeWithSelectAccount(ctx context.Context, data dto.AuthorizeContext) (string, error)
}

type switcher struct {
	login         LoginFlowProcessor
	none          NoneFlowProcessor
	consent       ConsentFlowProcessor
	selectAccount SelectAccountFlowProcessor
}

// NewFlowSwitcher returns new authorize switcher.
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

// Switch executes the authorize switcher.
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
