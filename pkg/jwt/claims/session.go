package claims

// SessionClaims contains the claims required to identify a user's session and
// manage its lifecycle.
//
// This structure complies with the OpenID Connect Front-Channel Logout 1.0
// specification (Section 3) and is used to coordinate single sign-out processes
// between the authorization server and client applications.
type SessionClaims struct {
	// SessionID is a unique identifier for a user's session.
	//
	// This token is critical for the Front-Channel and Back-Channel Logout
	// mechanisms. It associates the ID token with a specific browser session,
	// allowing the authorization server to send targeted session termination
	// requests to all trusted applications and services.
	SessionID string `json:"sid,omitempty"`
}

// NewSessionClaims returns a new SessionClaims struct with the provided unique
// session identifier.
func NewSessionClaims(sessionID string) SessionClaims {
	return SessionClaims{
		SessionID: sessionID,
	}
}
