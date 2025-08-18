package dto

type DataForLogin struct {
	User     User
	Client   Client
	Sessions []Session
}

type DataForRegister struct {
	User                         User
	Client                       Client
	ClientDefaultRoles           []Role
	ClientDefaultPermissionCodes []string
}

type DataForRefreshTokens struct {
	Session Session
	User    User
}

type DataForLogout struct {
	Session Session
}
