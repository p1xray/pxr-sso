package oauth

// Request parameter names
const (
	RequestParameterNameResponseType        = "response_type"
	RequestParameterNameClientID            = "client_id"
	RequestParameterNameRedirectURI         = "redirect_uri"
	RequestParameterNameCodeChallenge       = "code_challenge"
	RequestParameterNameCodeChallengeMethod = "code_challenge_method"
	RequestParameterNameState               = "state"
	RequestParameterNameFlowID              = "flow_id"
	RequestParameterNameUsername            = "username"
	RequestParameterNamePassword            = "password"
	RequestParameterNameAuthorizationCode   = "code"
	RequestParameterNameGrantType           = "grant_type"
	RequestParameterNameCodeVerifier        = "code_verifier"

	RequestParameterNameError            = "error"
	RequestParameterNameErrorDescription = "error_description"
	RequestParameterNameErrorURI         = "error_uri"
)

// Request parameter allowed values
const (
	RequestParameterAllowValueResponseTypeCode           = "code"
	RequestParameterAllowValueCodeChallengeMethod        = "S256"
	RequestParameterAllowValueGrantTypeAuthorizationCode = "authorization_code"
)

const (
	// RedisFlowTTL is a redis record TTL for flow (in minutes)
	RedisFlowTTL = 15

	// RedisAuthorizationCodeTTL is a redis record TTL for authorization code (in minutes)
	RedisAuthorizationCodeTTL = 10
)
