package dto

// DataForLogin is a DTO with data for logging in a user.
type DataForLogin struct {
	User     User
	Client   Client
	Sessions []Session
}

// DataForRegister is a DTO with data for registering a new user.
type DataForRegister struct {
	User                         User
	Client                       Client
	ClientDefaultRoles           []Role
	ClientDefaultPermissionCodes []string
}

// DataForRefreshTokens is a DTO with data for refreshing user tokens.
type DataForRefreshTokens struct {
	Session Session
	User    User
}

// DataForLogout is a DTO with data for logging out a user.
type DataForLogout struct {
	Session Session
}
