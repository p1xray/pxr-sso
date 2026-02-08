package authorize

type Params struct {
	ResponseType        []string
	ClientID            []string
	RedirectURI         []string
	CodeChallenge       []string
	CodeChallengeMethod []string
	State               []string
	Scope               []string
}
