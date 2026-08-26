package claims

// ClientClaims represents the client in OAuth 2.0 contexts.
//
// Defined in RFC8693 - OAuth 2.0 Token Exchange.
type ClientClaims struct {
	// ClientID represents the client identifier in OAuth 2.0 contexts.
	//
	// Identifies the OAuth 2.0 client that requested the token, providing a
	// mechanism for associating a token with a specific registered client
	// application, critical for enforcing client-specific access policies.
	ClientID string `json:"client_id,omitempty"`
}

// NewClientClaims returns a new ClientClaims struct with the specified OAuth 2.0
// client identifier.
func NewClientClaims(clientID string) ClientClaims {
	return ClientClaims{
		ClientID: clientID,
	}
}
