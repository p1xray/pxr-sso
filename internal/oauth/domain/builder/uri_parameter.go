package builder

import (
	"github.com/p1xray/pxr-sso/internal/oauth"
)

type queryParameters map[string]string

func newLoginQueryParameters(
	flowID,
	responseType,
	clientID,
	redirectURI,
	state,
	scope string,
) queryParameters {
	return queryParameters{
		oauth.RequestParameterNameFlowID:       flowID,
		oauth.RequestParameterNameResponseType: responseType,
		oauth.RequestParameterNameClientID:     clientID,
		oauth.RequestParameterNameRedirectURI:  redirectURI,
		oauth.RequestParameterNameState:        state,
		oauth.RequestParameterNameScope:        scope,
	}
}

func newConsentQueryParameters(
	flowID,
	responseType,
	clientID,
	redirectURI,
	state,
	scope string,
) queryParameters {
	return queryParameters{
		oauth.RequestParameterNameFlowID:       flowID,
		oauth.RequestParameterNameResponseType: responseType,
		oauth.RequestParameterNameClientID:     clientID,
		oauth.RequestParameterNameRedirectURI:  redirectURI,
		oauth.RequestParameterNameState:        state,
		oauth.RequestParameterNameScope:        scope,
	}
}

func newCallbackQueryParameters(authorizationCode, state string) queryParameters {
	return queryParameters{
		oauth.RequestParameterNameAuthorizationCode: authorizationCode,
		oauth.RequestParameterNameState:             state,
	}
}

func newErrorQueryParameters(code, description, errURI string) queryParameters {
	return queryParameters{
		oauth.RequestParameterNameError:            code,
		oauth.RequestParameterNameErrorDescription: description,
		oauth.RequestParameterNameErrorURI:         errURI,
	}
}
