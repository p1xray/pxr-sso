package claims

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const (
	testPhoneNumber         = "+1234567890"
	testPhoneNumberVerified = true
)

func TestNewOpenIDPhoneNumberClaims(t *testing.T) {
	expectedPhoneNumberVerified := testPhoneNumberVerified
	expectedClaims := OpenIDPhoneNumberClaims{
		PhoneNumber:         testPhoneNumber,
		PhoneNumberVerified: &expectedPhoneNumberVerified,
	}

	t.Run("successful creation of new OpenID phone claims", func(t *testing.T) {
		phoneClaims := NewOpenIDPhoneNumberClaims(testPhoneNumber, testPhoneNumberVerified)

		if !cmp.Equal(expectedClaims, phoneClaims) {
			t.Fatal(cmp.Diff(expectedClaims, phoneClaims))
		}
	})
}
