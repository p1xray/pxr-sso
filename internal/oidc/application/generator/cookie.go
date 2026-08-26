package generator

const sessionCookieNamePrefix = "pxr.sso.session_"

func SessionCookieName(code string) string {
	return sessionCookieNamePrefix + code
}
