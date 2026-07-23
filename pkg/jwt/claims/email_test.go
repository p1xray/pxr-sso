package claims

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const (
	testEmail         = "johndoe@example.com"
	testEmailVerified = false
)

func TestNewOpenIDEmailClaims(t *testing.T) {
	expectedEmailVerified := testEmailVerified
	expectedClaims := OpenIDEmailClaims{
		Email:         testEmail,
		EmailVerified: &expectedEmailVerified,
	}

	t.Run("successful creation of new OpenID email claims", func(t *testing.T) {
		emailClaims := NewOpenIDEmailClaims(testEmail, testEmailVerified)

		if !cmp.Equal(expectedClaims, emailClaims) {
			t.Fatal(cmp.Diff(expectedClaims, emailClaims))
		}
	})
}
