package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestNewIDTokenClaims(t *testing.T) {
	expectedPhoneNumberVerified := testPhoneNumberVerified
	expectedEmailVerified := testEmailVerified

	testCases := []struct {
		name           string
		issuer         string
		subject        string
		audience       jwt.Audience
		ttl            time.Duration
		opts           []IDTokenClaimsOption
		expectedClaims IDTokenClaims
	}{
		{
			name:     "successful creation of new id token claims with required parameters",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      testTTL,
			expectedClaims: IDTokenClaims{
				RegisteredClaims: RegisteredClaims{
					Issuer:   testIssuer,
					Subject:  testSubject,
					Audience: testAudience,
				},
			},
		},
		{
			name:     "successful creation of new id token claims with all required and optional parameters",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      testTTL,
			opts: []IDTokenClaimsOption{
				WithIDTokenAuthentication(
					testAuthTime,
					WithNonce(testNonce),
					WithReferences(testACR, testAMR),
					WithAuthorizedParty(testAZP),
					WithHash(testAtHash, testCHash),
				),
				WithIDTokenProfile(
					WithProfileName(testName, testFamilyName, testGivenName, testMiddleName),
					WithProfileNameAlias(testNickname, testPreferredUsername),
					WithProfileURL(testProfile, testPicture, testWebsite),
					WithProfileGender(testGender),
					WithProfileBirthdate(testBirthdate),
					WithProfileLocation(testZoneInfo, testLocale),
					WithProfileUpdatedAt(updatedAt),
				),
				WithIDTokenPhoneNumber(testPhoneNumber, testPhoneNumberVerified),
				WithIDTokenEmail(testEmail, testEmailVerified),
				WithIDTokenAddress(
					testAddressFormatted,
					testStreetAddress,
					testLocality,
					testRegion,
					testPostalCode,
					testCountry,
				),
			},
			expectedClaims: IDTokenClaims{
				RegisteredClaims: RegisteredClaims{
					Issuer:   testIssuer,
					Subject:  testSubject,
					Audience: testAudience,
				},
				AuthenticationClaims: AuthenticationClaims{
					AuthTime: jwt.NewNumericDate(testAuthTime),
					Nonce:    testNonce,
					ACR:      testACR,
					AMR:      testAMR,
					AZP:      testAZP,
					AtHash:   testAtHash,
					CHash:    testCHash,
				},
				OpenIDProfileClaims: OpenIDProfileClaims{
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
				},
				OpenIDPhoneNumberClaims: OpenIDPhoneNumberClaims{
					PhoneNumber:         testPhoneNumber,
					PhoneNumberVerified: &expectedPhoneNumberVerified,
				},
				OpenIDEmailClaims: OpenIDEmailClaims{
					Email:         testEmail,
					EmailVerified: &expectedEmailVerified,
				},
				OpenIDAddressClaims: OpenIDAddressClaims{
					Address: &AddressClaim{
						Formatted:     testAddressFormatted,
						StreetAddress: testStreetAddress,
						Locality:      testLocality,
						Region:        testRegion,
						PostalCode:    testPostalCode,
						Country:       testCountry,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenClaims := NewIDTokenClaims(tc.issuer, tc.subject, tc.audience, tc.ttl, tc.opts...)

			checkRegisteredClaims(
				t,
				tokenClaims.RegisteredClaims,
				tc.expectedClaims.Issuer,
				tc.expectedClaims.Subject,
				tc.expectedClaims.Audience,
				leeway,
			)

			if !cmp.Equal(tc.expectedClaims.AuthenticationClaims, tokenClaims.AuthenticationClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.AuthenticationClaims, tokenClaims.AuthenticationClaims))
			}

			if !cmp.Equal(tc.expectedClaims.OpenIDProfileClaims, tokenClaims.OpenIDProfileClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.OpenIDProfileClaims, tokenClaims.OpenIDProfileClaims))
			}

			if !cmp.Equal(tc.expectedClaims.OpenIDPhoneNumberClaims, tokenClaims.OpenIDPhoneNumberClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.OpenIDPhoneNumberClaims, tokenClaims.OpenIDPhoneNumberClaims))
			}

			if !cmp.Equal(tc.expectedClaims.OpenIDEmailClaims, tokenClaims.OpenIDEmailClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.OpenIDEmailClaims, tokenClaims.OpenIDEmailClaims))
			}

			if !cmp.Equal(tc.expectedClaims.OpenIDAddressClaims, tokenClaims.OpenIDAddressClaims) {
				t.Fatal(cmp.Diff(tc.expectedClaims.OpenIDAddressClaims, tokenClaims.OpenIDAddressClaims))
			}
		})
	}
}
