package claims

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const testClientID = "test_client_id"

func TestNewClientClaims(t *testing.T) {
	expectedClaims := ClientClaims{
		ClientID: testClientID,
	}

	t.Run("successful creation of new client claims", func(t *testing.T) {
		clientClaims := NewClientClaims(testClientID)

		if !cmp.Equal(expectedClaims, clientClaims) {
			t.Fatal(cmp.Diff(expectedClaims, clientClaims))
		}
	})
}
