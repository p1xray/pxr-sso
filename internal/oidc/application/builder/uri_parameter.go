package builder

import (
	"github.com/p1xray/pxr-sso/internal/oidc"
)

type queryParameters map[string]string

func newRequestURIQueryParameters(requestURI string) queryParameters {
	return queryParameters{
		oidc.RequestParameterNameRequestURI: requestURI,
	}
}

func newCallbackQueryParameters(authorizationCode, state string) queryParameters {
	return queryParameters{
		oidc.RequestParameterNameAuthorizationCode: authorizationCode,
		oidc.RequestParameterNameState:             state,
	}
}

func newErrorQueryParameters(code, description, errURI string) queryParameters {
	return queryParameters{
		oidc.RequestParameterNameError:            code,
		oidc.RequestParameterNameErrorDescription: description,
		oidc.RequestParameterNameErrorURI:         errURI,
	}
}
