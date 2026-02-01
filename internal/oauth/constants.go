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

	RequestParameterNameError            = "error"
	RequestParameterNameErrorDescription = "error_description"
	RequestParameterNameErrorURI         = "error_uri"
)

// Request parameter allowed values
const (
	RequestParameterAllowValueResponseTypeCode    = "code"
	RequestParameterAllowValueCodeChallengeMethod = "S256"
)

// Redis object type names
const (
	RedisObjectTypeNameFlow = "flow"
)
