package claims

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const (
	testAddressFormatted = "123 Main St\nMetropolis, NY 10001\nUSA"
	testStreetAddress    = "123 Main St"
	testLocality         = "Metropolis"
	testRegion           = "NY"
	testPostalCode       = "10001"
	testCountry          = "USA"
)

func TestNewOpenIDAddressClaims(t *testing.T) {
	expectedClaims := OpenIDAddressClaims{
		Address: &AddressClaim{
			Formatted:     testAddressFormatted,
			StreetAddress: testStreetAddress,
			Locality:      testLocality,
			Region:        testRegion,
			PostalCode:    testPostalCode,
			Country:       testCountry,
		},
	}

	t.Run("successful creation of new OpenID address claims", func(t *testing.T) {
		addressClaims := NewOpenIDAddressClaims(
			testAddressFormatted,
			testStreetAddress,
			testLocality,
			testRegion,
			testPostalCode,
			testCountry,
		)

		if !cmp.Equal(expectedClaims, addressClaims) {
			t.Fatal(cmp.Diff(expectedClaims, addressClaims))
		}
	})
}
