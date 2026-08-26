package jwt

import (
	"errors"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

func TestCreateAccessToken(t *testing.T) {
	testCases := []struct {
		name          string
		claims        claims.AccessTokenClaims
		key           []byte
		expectedError error
	}{
		{
			name:   "successfully creates a new access token with minimal set of claims",
			claims: NewTestMinimalAccessTokenClaims(),
			key:    []byte(testValidKey),
		},
		{
			name:   "successfully creates a new access token with full set of claims",
			claims: NewTestFullAccessTokenClaims(),
			key:    []byte(testValidKey),
		},
		{
			name:          "throws an error when creating a token signed by invalid key",
			claims:        NewTestMinimalAccessTokenClaims(),
			key:           []byte(testInvalidKey),
			expectedError: ErrTokenSerialize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rawToken, err := CreateAccessToken(tc.claims, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				token, err := jwt.ParseSigned(rawToken, []jose.SignatureAlgorithm{jose.HS256})
				require.NoError(t, err)

				tokenClaims := make(map[string]interface{})
				err = token.Claims(tc.key, &tokenClaims)
				require.NoError(t, err)

				checkAccessTokenClaims(t, tokenClaims, tc.claims)
			}
		})
	}
}

func TestCreateRefreshToken(t *testing.T) {
	testCases := []struct {
		name          string
		claims        claims.RefreshTokenClaims
		key           []byte
		expectedError error
	}{
		{
			name:   "successfully creates a new refresh token with minimal set of claims",
			claims: NewTestMinimalRefreshTokenClaims(),
			key:    []byte(testValidKey),
		},
		{
			name:   "successfully creates a new refresh token with full set of claims",
			claims: NewTestFullRefreshTokenClaims(),
			key:    []byte(testValidKey),
		},
		{
			name:          "throws an error when creating a token signed by invalid key",
			claims:        NewTestMinimalRefreshTokenClaims(),
			key:           []byte(testInvalidKey),
			expectedError: ErrTokenSerialize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rawToken, err := CreateRefreshToken(tc.claims, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				token, err := jwt.ParseSigned(rawToken, []jose.SignatureAlgorithm{jose.HS256})
				require.NoError(t, err)

				tokenClaims := make(map[string]interface{})
				err = token.Claims(tc.key, &tokenClaims)
				require.NoError(t, err)

				checkRefreshTokenClaims(t, tokenClaims, tc.claims)
			}
		})
	}
}

func TestCreateIDToken(t *testing.T) {
	testCases := []struct {
		name          string
		claims        claims.IDTokenClaims
		key           []byte
		expectedError error
	}{
		{
			name:   "successfully creates a new id token with minimal set of claims",
			claims: NewTestMinimalIDTokenClaims(),
			key:    []byte(testValidKey),
		},
		{
			name:   "successfully creates a new id token with full set of claims",
			claims: NewTestFullIDTokenClaims(),
			key:    []byte(testValidKey),
		},
		{
			name:          "throws an error when creating a token signed by invalid key",
			claims:        NewTestMinimalIDTokenClaims(),
			key:           []byte(testInvalidKey),
			expectedError: ErrTokenSerialize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rawToken, err := CreateIDToken(tc.claims, tc.key)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				token, err := jwt.ParseSigned(rawToken, []jose.SignatureAlgorithm{jose.HS256})
				require.NoError(t, err)

				tokenClaims := make(map[string]interface{})
				err = token.Claims(tc.key, &tokenClaims)
				require.NoError(t, err)

				checkIDTokenClaims(t, tokenClaims, tc.claims)
			}
		})
	}
}

func checkAccessTokenClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.AccessTokenClaims) {
	checkRegisteredClaims(t, tokenClaims, expectedClaims.RegisteredClaims)
	checkAuthenticationClaims(t, tokenClaims, expectedClaims.AuthenticationClaims)
	checkAuthorizationClaims(t, tokenClaims, expectedClaims.AuthorizationClaims)
	checkClientClaims(t, tokenClaims, expectedClaims.ClientClaims)
}

func checkRefreshTokenClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.RefreshTokenClaims) {
	checkRegisteredClaims(t, tokenClaims, expectedClaims.RegisteredClaims)
	checkAuthenticationClaims(t, tokenClaims, expectedClaims.AuthenticationClaims)
	checkAuthorizationClaims(t, tokenClaims, expectedClaims.AuthorizationClaims)
	checkClientClaims(t, tokenClaims, expectedClaims.ClientClaims)
	checkSessionClaims(t, tokenClaims, expectedClaims.SessionClaims)
}

func checkIDTokenClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.IDTokenClaims) {
	checkRegisteredClaims(t, tokenClaims, expectedClaims.RegisteredClaims)
	checkAuthenticationClaims(t, tokenClaims, expectedClaims.AuthenticationClaims)
	checkOpenIDProfileClaims(t, tokenClaims, expectedClaims.OpenIDProfileClaims)
	checkOpenIDPhoneNumberClaims(t, tokenClaims, expectedClaims.OpenIDPhoneNumberClaims)
	checkOpenIDEmailClaims(t, tokenClaims, expectedClaims.OpenIDEmailClaims)
	checkOpenIDAddressClaims(t, tokenClaims, expectedClaims.OpenIDAddressClaims)
}

func checkRegisteredClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.RegisteredClaims) {
	checkJtiClaim(t, tokenClaims, expectedClaims.ID)
	checkSubClaim(t, tokenClaims, expectedClaims.Subject)
	checkIssClaim(t, tokenClaims, expectedClaims.Issuer)
	checkAudClaim(t, tokenClaims, expectedClaims.Audience)
	checkExpClaim(t, tokenClaims, fixedDate)
	checkNbfClaim(t, tokenClaims, fixedDate)
	checkIatClaim(t, tokenClaims, fixedDate)
}

func checkJtiClaim(t *testing.T, tokenClaims map[string]interface{}, expectedID string) {
	jti, ok := tokenClaims["jti"]
	require.True(t, ok)

	assert.Equal(t, expectedID, jti.(string))
}

func checkSubClaim(t *testing.T, tokenClaims map[string]interface{}, expectedSubject string) {
	sub, ok := tokenClaims["sub"]
	require.True(t, ok)

	assert.Equal(t, expectedSubject, sub.(string))
}

func checkIssClaim(t *testing.T, tokenClaims map[string]interface{}, expectedIssuer string) {
	iss, ok := tokenClaims["iss"]
	require.True(t, ok)

	assert.Equal(t, expectedIssuer, iss.(string))
}

func checkAudClaim(t *testing.T, tokenClaims map[string]interface{}, expectedAudiences []string) {
	aud, ok := tokenClaims["aud"]
	require.True(t, ok)

	audiences, err := convertClaimValueToStringArray(aud)
	require.NoError(t, err)

	assert.Equal(t, expectedAudiences, audiences)
}

func checkExpClaim(t *testing.T, tokenClaims map[string]interface{}, date time.Time) {
	exp, ok := tokenClaims["exp"]
	require.True(t, ok)

	expTime := time.Unix(int64(exp.(float64)), 0).UTC()
	expValid := expTime.After(date)
	assert.True(t, expValid)
}

func checkNbfClaim(t *testing.T, tokenClaims map[string]interface{}, date time.Time) {
	nbf, ok := tokenClaims["nbf"]
	require.True(t, ok)

	nbfTime := time.Unix(int64(nbf.(float64)), 0).UTC()
	nbfValid := nbfTime.Before(date)
	assert.True(t, nbfValid)
}

func checkIatClaim(t *testing.T, tokenClaims map[string]interface{}, date time.Time) {
	iat, ok := tokenClaims["iat"]
	require.True(t, ok)

	iatTime := time.Unix(int64(iat.(float64)), 0).UTC()
	iatValid := iatTime.Before(date)
	assert.True(t, iatValid)
}

func checkAuthenticationClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.AuthenticationClaims) {
	checkAuthTimeClaim(t, tokenClaims, expectedClaims.AuthTime)
	checkNonceClaim(t, tokenClaims, expectedClaims.Nonce)
	checkACRClaim(t, tokenClaims, expectedClaims.ACR)
	checkAMRClaim(t, tokenClaims, expectedClaims.AMR)
	checkAZPClaim(t, tokenClaims, expectedClaims.AZP)
	checkAtHashClaim(t, tokenClaims, expectedClaims.AtHash)
	checkCHashClaim(t, tokenClaims, expectedClaims.CHash)
}

func checkAuthTimeClaim(t *testing.T, tokenClaims map[string]interface{}, expectedAuthTime *jwt.NumericDate) {
	if expectedAuthTime == nil {
		return
	}

	authTimeClaim, ok := tokenClaims["auth_time"]
	require.True(t, ok)

	authTime := int64(authTimeClaim.(float64))

	assert.Equal(t, claims.NumericDateToInt64(expectedAuthTime), authTime)
}

func checkNonceClaim(t *testing.T, tokenClaims map[string]interface{}, expectedNonce string) {
	if expectedNonce == "" {
		return
	}

	nonce, ok := tokenClaims["nonce"]
	require.True(t, ok)

	assert.Equal(t, expectedNonce, nonce.(string))
}

func checkACRClaim(t *testing.T, tokenClaims map[string]interface{}, expectedACR string) {
	if expectedACR == "" {
		return
	}

	acr, ok := tokenClaims["acr"]
	require.True(t, ok)

	assert.Equal(t, expectedACR, acr.(string))
}

func checkAMRClaim(t *testing.T, tokenClaims map[string]interface{}, expectedAMR []string) {
	if expectedAMR == nil {
		return
	}

	amr, ok := tokenClaims["amr"]
	require.True(t, ok)

	references, err := convertClaimValueToStringArray(amr)
	require.NoError(t, err)

	assert.Equal(t, expectedAMR, references)
}

func checkAZPClaim(t *testing.T, tokenClaims map[string]interface{}, expectedAZP string) {
	if expectedAZP == "" {
		return
	}

	azp, ok := tokenClaims["azp"]
	require.True(t, ok)

	assert.Equal(t, expectedAZP, azp.(string))
}

func checkAtHashClaim(t *testing.T, tokenClaims map[string]interface{}, expectedAtHash string) {
	if expectedAtHash == "" {
		return
	}

	atHash, ok := tokenClaims["at_hash"]
	require.True(t, ok)

	assert.Equal(t, expectedAtHash, atHash.(string))
}

func checkCHashClaim(t *testing.T, tokenClaims map[string]interface{}, expectedCHash string) {
	if expectedCHash == "" {
		return
	}

	cHash, ok := tokenClaims["c_hash"]
	require.True(t, ok)

	assert.Equal(t, expectedCHash, cHash.(string))
}

func checkAuthorizationClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.AuthorizationClaims) {
	checkScopeClaim(t, tokenClaims, expectedClaims.Scope)
	checkRolesClaim(t, tokenClaims, expectedClaims.Roles)
	checkGroupsClaim(t, tokenClaims, expectedClaims.Groups)
	checkEntitlementsClaim(t, tokenClaims, expectedClaims.Entitlements)
}

func checkScopeClaim(t *testing.T, tokenClaims map[string]interface{}, expectedScope string) {
	if expectedScope == "" {
		return
	}

	scope, ok := tokenClaims["scope"]
	require.True(t, ok)

	assert.Equal(t, expectedScope, scope.(string))
}

func checkRolesClaim(t *testing.T, tokenClaims map[string]interface{}, expectedRoles []string) {
	if expectedRoles == nil || len(expectedRoles) == 0 {
		return
	}

	rolesClaim, ok := tokenClaims["roles"]
	require.True(t, ok)

	roles, err := convertClaimValueToStringArray(rolesClaim)
	require.NoError(t, err)

	assert.Equal(t, expectedRoles, roles)
}

func checkGroupsClaim(t *testing.T, tokenClaims map[string]interface{}, expectedGroups []string) {
	if expectedGroups == nil || len(expectedGroups) == 0 {
		return
	}

	groupsClaim, ok := tokenClaims["groups"]
	require.True(t, ok)

	groups, err := convertClaimValueToStringArray(groupsClaim)
	require.NoError(t, err)

	assert.Equal(t, expectedGroups, groups)
}

func checkEntitlementsClaim(t *testing.T, tokenClaims map[string]interface{}, expectedEntitlements []string) {
	if expectedEntitlements == nil || len(expectedEntitlements) == 0 {
		return
	}

	entitlementsClaim, ok := tokenClaims["entitlements"]
	require.True(t, ok)

	entitlements, err := convertClaimValueToStringArray(entitlementsClaim)
	require.NoError(t, err)

	assert.Equal(t, expectedEntitlements, entitlements)
}

func checkClientClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.ClientClaims) {
	checkClientIDClaim(t, tokenClaims, expectedClaims.ClientID)
}

func checkClientIDClaim(t *testing.T, tokenClaims map[string]interface{}, expectedClientID string) {
	if expectedClientID == "" {
		return
	}

	clientClaim, ok := tokenClaims["client_id"]
	require.True(t, ok)

	assert.Equal(t, expectedClientID, clientClaim.(string))
}

func checkSessionClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.SessionClaims) {
	checkSessionIDClaim(t, tokenClaims, expectedClaims.SessionID)
}

func checkSessionIDClaim(t *testing.T, tokenClaims map[string]interface{}, expectedSessionID string) {
	if expectedSessionID == "" {
		return
	}

	sid, ok := tokenClaims["sid"]
	require.True(t, ok)

	assert.Equal(t, expectedSessionID, sid.(string))
}

func checkOpenIDProfileClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.OpenIDProfileClaims) {
	checkNameClaim(t, tokenClaims, expectedClaims.Name)
	checkFamilyNameClaim(t, tokenClaims, expectedClaims.FamilyName)
	checkGivenNameClaim(t, tokenClaims, expectedClaims.GivenName)
	checkMiddleNameClaim(t, tokenClaims, expectedClaims.MiddleName)
	checkNicknameClaim(t, tokenClaims, expectedClaims.Nickname)
	checkPreferredUsernameClaim(t, tokenClaims, expectedClaims.PreferredUsername)
	checkProfileClaim(t, tokenClaims, expectedClaims.Profile)
	checkPictureClaim(t, tokenClaims, expectedClaims.Picture)
	checkWebsiteClaim(t, tokenClaims, expectedClaims.Website)
	checkGenderClaim(t, tokenClaims, expectedClaims.Gender)
	checkBirthdateClaim(t, tokenClaims, expectedClaims.Birthdate)
	checkZoneInfoClaim(t, tokenClaims, expectedClaims.ZoneInfo)
	checkLocaleClaim(t, tokenClaims, expectedClaims.Locale)
	checkUpdatedAtClaim(t, tokenClaims, expectedClaims.UpdatedAt)
}

func checkNameClaim(t *testing.T, tokenClaims map[string]interface{}, expectedName string) {
	if expectedName == "" {
		return
	}

	name, ok := tokenClaims["name"]
	require.True(t, ok)

	assert.Equal(t, expectedName, name.(string))
}

func checkFamilyNameClaim(t *testing.T, tokenClaims map[string]interface{}, expectedFamilyName string) {
	if expectedFamilyName == "" {
		return
	}

	familyName, ok := tokenClaims["family_name"]
	require.True(t, ok)

	assert.Equal(t, expectedFamilyName, familyName.(string))
}

func checkGivenNameClaim(t *testing.T, tokenClaims map[string]interface{}, expectedGivenName string) {
	if expectedGivenName == "" {
		return
	}

	givenName, ok := tokenClaims["given_name"]
	require.True(t, ok)

	assert.Equal(t, expectedGivenName, givenName.(string))
}

func checkMiddleNameClaim(t *testing.T, tokenClaims map[string]interface{}, expectedMiddleName string) {
	if expectedMiddleName == "" {
		return
	}

	middleName, ok := tokenClaims["middle_name"]
	require.True(t, ok)

	assert.Equal(t, expectedMiddleName, middleName.(string))
}

func checkNicknameClaim(t *testing.T, tokenClaims map[string]interface{}, expectedNickname string) {
	if expectedNickname == "" {
		return
	}

	nickname, ok := tokenClaims["nickname"]
	require.True(t, ok)

	assert.Equal(t, expectedNickname, nickname.(string))
}

func checkPreferredUsernameClaim(t *testing.T, tokenClaims map[string]interface{}, expectedPreferredUsername string) {
	if expectedPreferredUsername == "" {
		return
	}

	preferredUsername, ok := tokenClaims["preferred_username"]
	require.True(t, ok)

	assert.Equal(t, expectedPreferredUsername, preferredUsername.(string))
}

func checkProfileClaim(t *testing.T, tokenClaims map[string]interface{}, expectedProfile string) {
	if expectedProfile == "" {
		return
	}

	profile, ok := tokenClaims["profile"]
	require.True(t, ok)

	_, err := url.Parse(profile.(string))
	assert.NoError(t, err)

	assert.Equal(t, expectedProfile, profile.(string))
}

func checkPictureClaim(t *testing.T, tokenClaims map[string]interface{}, expectedPicture string) {
	if expectedPicture == "" {
		return
	}

	picture, ok := tokenClaims["picture"]
	require.True(t, ok)

	_, err := url.Parse(picture.(string))
	assert.NoError(t, err)

	assert.Equal(t, expectedPicture, picture.(string))
}

func checkWebsiteClaim(t *testing.T, tokenClaims map[string]interface{}, expectedWebsite string) {
	if expectedWebsite == "" {
		return
	}

	website, ok := tokenClaims["website"]
	require.True(t, ok)

	_, err := url.Parse(website.(string))
	assert.NoError(t, err)

	assert.Equal(t, expectedWebsite, website.(string))
}

func checkGenderClaim(t *testing.T, tokenClaims map[string]interface{}, expectedGender string) {
	if expectedGender == "" {
		return
	}

	gender, ok := tokenClaims["gender"]
	require.True(t, ok)

	assert.Equal(t, expectedGender, gender.(string))
}

func checkBirthdateClaim(t *testing.T, tokenClaims map[string]interface{}, expectedBirthdate string) {
	if expectedBirthdate == "" {
		return
	}

	birthdate, ok := tokenClaims["birthdate"]
	require.True(t, ok)

	valid := isValidDate(birthdate.(string))
	assert.True(t, valid)

	assert.Equal(t, expectedBirthdate, birthdate.(string))
}

func checkZoneInfoClaim(t *testing.T, tokenClaims map[string]interface{}, expectedZoneInfo string) {
	if expectedZoneInfo == "" {
		return
	}

	zoneInfo, ok := tokenClaims["zoneinfo"]
	require.True(t, ok)

	assert.Equal(t, expectedZoneInfo, zoneInfo.(string))
}

func checkLocaleClaim(t *testing.T, tokenClaims map[string]interface{}, expectedLocale string) {
	if expectedLocale == "" {
		return
	}

	locale, ok := tokenClaims["locale"]
	require.True(t, ok)

	assert.Equal(t, expectedLocale, locale.(string))
}

func checkUpdatedAtClaim(t *testing.T, tokenClaims map[string]interface{}, expectedUpdatedAt *jwt.NumericDate) {
	if expectedUpdatedAt == nil {
		return
	}

	updatedAtClaim, ok := tokenClaims["updated_at"]
	require.True(t, ok)

	updatedAt := int64(updatedAtClaim.(float64))

	assert.Equal(t, claims.NumericDateToInt64(expectedUpdatedAt), updatedAt)
}

func checkOpenIDPhoneNumberClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.OpenIDPhoneNumberClaims) {
	checkPhoneNumberClaim(t, tokenClaims, expectedClaims.PhoneNumber)
	checkPhoneNumberVerifiedClaim(t, tokenClaims, expectedClaims.PhoneNumberVerified)
}

func checkPhoneNumberClaim(t *testing.T, tokenClaims map[string]interface{}, expectedPhoneNumber string) {
	if expectedPhoneNumber == "" {
		return
	}

	phoneNumber, ok := tokenClaims["phone_number"]
	require.True(t, ok)

	assert.Equal(t, expectedPhoneNumber, phoneNumber.(string))
}

func checkPhoneNumberVerifiedClaim(t *testing.T, tokenClaims map[string]interface{}, expectedPhoneNumberVerified *bool) {
	if expectedPhoneNumberVerified == nil {
		return
	}

	phoneNumberVerified, ok := tokenClaims["phone_number_verified"]
	require.True(t, ok)

	assert.Equal(t, *expectedPhoneNumberVerified, phoneNumberVerified.(bool))
}

func checkOpenIDEmailClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.OpenIDEmailClaims) {
	checkEmailClaim(t, tokenClaims, expectedClaims.Email)
	checkEmailVerifiedClaim(t, tokenClaims, expectedClaims.EmailVerified)

}

func checkEmailClaim(t *testing.T, tokenClaims map[string]interface{}, expectedEmail string) {
	if expectedEmail == "" {
		return
	}

	email, ok := tokenClaims["email"]
	require.True(t, ok)

	assert.Equal(t, expectedEmail, email.(string))
}

func checkEmailVerifiedClaim(t *testing.T, tokenClaims map[string]interface{}, expectedEmailVerified *bool) {
	if expectedEmailVerified == nil {
		return
	}

	emailVerified, ok := tokenClaims["email_verified"]
	require.True(t, ok)

	assert.Equal(t, *expectedEmailVerified, emailVerified.(bool))
}

func checkOpenIDAddressClaims(t *testing.T, tokenClaims map[string]interface{}, expectedClaims claims.OpenIDAddressClaims) {
	checkAddressClaim(t, tokenClaims, expectedClaims.Address)
}

func checkAddressClaim(t *testing.T, tokenClaims map[string]interface{}, expectedAddress *claims.AddressClaim) {
	if expectedAddress == nil {
		return
	}

	rawAddress, ok := tokenClaims["address"]
	require.True(t, ok)

	addressClaims, ok := rawAddress.(map[string]interface{})
	require.True(t, ok)

	checkFormattedClaim(t, addressClaims, expectedAddress.Formatted)
	checkStreetAddressClaim(t, addressClaims, expectedAddress.StreetAddress)
	checkLocalityClaim(t, addressClaims, expectedAddress.Locality)
	checkRegionClaim(t, addressClaims, expectedAddress.Region)
	checkPostalCodeClaim(t, addressClaims, expectedAddress.PostalCode)
	checkCountryClaim(t, addressClaims, expectedAddress.Country)
}

func checkFormattedClaim(t *testing.T, addressMap map[string]interface{}, expected string) {
	val, ok := addressMap["formatted"]
	require.True(t, ok)

	assert.Equal(t, expected, val.(string))
}

func checkStreetAddressClaim(t *testing.T, addressMap map[string]interface{}, expected string) {
	val, ok := addressMap["street_address"]
	require.True(t, ok)

	assert.Equal(t, expected, val.(string))
}

func checkLocalityClaim(t *testing.T, addressMap map[string]interface{}, expected string) {
	val, ok := addressMap["locality"]
	require.True(t, ok)

	assert.Equal(t, expected, val.(string))
}

func checkRegionClaim(t *testing.T, addressMap map[string]interface{}, expected string) {
	val, ok := addressMap["region"]
	require.True(t, ok)

	assert.Equal(t, expected, val.(string))
}

func checkPostalCodeClaim(t *testing.T, addressMap map[string]interface{}, expected string) {
	val, ok := addressMap["postal_code"]
	require.True(t, ok)

	assert.Equal(t, expected, val.(string))
}

func checkCountryClaim(t *testing.T, addressMap map[string]interface{}, expected string) {
	val, ok := addressMap["country"]
	require.True(t, ok)

	assert.Equal(t, expected, val.(string))
}

func convertClaimValueToStringArray(claim any) ([]string, error) {
	result := make([]string, 0)
	switch claimValue := claim.(type) {
	case string:
		result = []string{claimValue}
	case []any:
		result = make([]string, len(claimValue))
		for i, v := range claimValue {
			audValueStr, ok := v.(string)
			if !ok {
				return nil, errors.New("invalid type of claim value into array")
			}

			result[i] = audValueStr
		}
	default:
		return nil, errors.New("invalid type of claim")
	}

	return result, nil
}

func isValidDate(dateStr string) bool {
	const layout = "2006-01-02"
	_, err := time.Parse(layout, dateStr)
	return err == nil
}
