package claims

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const (
	testSessionID = "test_session_id"
)

func TestNewSessionClaims(t *testing.T) {
	expectedClaims := SessionClaims{
		SessionID: testSessionID,
	}

	t.Run("successful creation of new session claims", func(t *testing.T) {
		sessionClaims := NewSessionClaims(testSessionID)

		if !cmp.Equal(expectedClaims, sessionClaims) {
			t.Fatal(cmp.Diff(expectedClaims, sessionClaims))
		}
	})
}
