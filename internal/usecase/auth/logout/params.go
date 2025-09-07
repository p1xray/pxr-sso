package logout

// Params is a data for log out use-case.
type Params struct {
	RefreshToken string
	ClientCode   string
}
