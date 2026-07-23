package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

const (
	testName              = "John Michael Doe"
	testFamilyName        = "Doe"
	testGivenName         = "John"
	testMiddleName        = "Michael"
	testNickname          = "Johnny"
	testPreferredUsername = "johndoe"
	testProfile           = "https://example.com/profile/johndoe"
	testPicture           = "https://example.com/profile/picture/johndoe"
	testWebsite           = "https://johndoe.dev"
	testGender            = "male"
	testBirthdate         = "1990-01-01"
	testZoneInfo          = "America/New_York"
	testLocale            = "en-US"
)

var (
	updatedAt = time.Date(2026, time.July, 21, 12, 0, 0, 0, time.UTC)
)

func TestNewOpenIDProfileClaims(t *testing.T) {
	expectedClaims := OpenIDProfileClaims{
		Name:              testName,
		FamilyName:        testFamilyName,
		GivenName:         testGivenName,
		MiddleName:        testMiddleName,
		Nickname:          testNickname,
		PreferredUsername: testPreferredUsername,
		Profile:           testProfile,
		Picture:           testPicture,
		Website:           testWebsite,
		Gender:            testGender,
		Birthdate:         testBirthdate,
		ZoneInfo:          testZoneInfo,
		Locale:            testLocale,
		UpdatedAt:         jwt.NewNumericDate(updatedAt),
	}

	t.Run("successful creation of new OpenID profile claims with all optional parameters", func(t *testing.T) {
		profileClaims := NewOpenIDProfileClaims(
			WithProfileName(testName, testFamilyName, testGivenName, testMiddleName),
			WithProfileNameAlias(testNickname, testPreferredUsername),
			WithProfileURL(testProfile, testPicture, testWebsite),
			WithProfileGender(testGender),
			WithProfileBirthdate(testBirthdate),
			WithProfileLocation(testZoneInfo, testLocale),
			WithProfileUpdatedAt(updatedAt),
		)

		if !cmp.Equal(expectedClaims, profileClaims) {
			t.Fatal(cmp.Diff(expectedClaims, profileClaims))
		}
	})
}
