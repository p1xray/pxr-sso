package generator

const sessionCookieNamePrefix = "pxr.sso.session:"

func SessionCookieName(code string) string {
	return sessionCookieNamePrefix + code
}
