package validator

// Request parameter allowed values.
const (
	RequestParameterAllowValueResponseTypeCode           = "code"
	RequestParameterAllowValueCodeChallengeMethod        = "S256"
	RequestParameterAllowValueGrantTypeAuthorizationCode = "authorization_code"
	RequestParameterAllowValuePromptNone                 = "none"
	RequestParameterAllowValuePromptLogin                = "login"
	RequestParameterAllowValuePromptConsent              = "consent"
	RequestParameterAllowValuePromptSelectAccount        = "select_account"
)

var RequestParameterAllowValuePrompt = []string{
	RequestParameterAllowValuePromptNone,
	RequestParameterAllowValuePromptLogin,
	RequestParameterAllowValuePromptConsent,
	RequestParameterAllowValuePromptSelectAccount,
}
