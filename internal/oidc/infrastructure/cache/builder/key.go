package builder

// authorizationRequestKeyPrefix is the prefix for a storage key for an authorization request.
const authorizationRequestKeyPrefix = "pxr.sso:par:"

// authorizedGrantKeyPrefix is the prefix for a storage key for an authorized grant.
const authorizedGrantKeyPrefix = "pxr.sso:grant:"

// AuthorizationRequestKey generates a storage key for an authorization request by URI.
func AuthorizationRequestKey(requestURI string) string {
	return authorizationRequestKeyPrefix + requestURI
}

// AuthorizedGrantKey generates a storage key for an authorized grant by authorization code.
func AuthorizedGrantKey(code string) string {
	return authorizedGrantKeyPrefix + code
}
