package builder

// authorizationRequestKeyPrefix is the prefix for a storage key for an authorization request.
const authorizationRequestKeyPrefix = "pxr.sso:par:"

// AuthorizationRequestKey generates a storage key for an authorization request by URI.
func AuthorizationRequestKey(requestURI string) string {
	return authorizationRequestKeyPrefix + requestURI
}
