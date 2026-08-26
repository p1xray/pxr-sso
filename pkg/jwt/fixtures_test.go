package jwt

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/p1xray/pxr-sso/pkg/jwt/claims"
	"time"
)

const (
	testValidKey   = "05c5328f-17cb-4b42-a085-4089c03b86f8"
	testInvalidKey = "invalid_key"

	testAccessTokenID  = "a0c43ab3-3067-4211-8086-f6c7f3009a82"
	testRefreshTokenID = "ef6b5f60-6700-4d07-a0a7-9b5604305304"
	testIDTokenID      = "3eeda4bb-96d1-44b6-8ed2-d59757c786bd"
	testIssuer         = "https://example.com/"
	testAudience       = "testAudience"
	testSubject        = "1"
	testScope          = "openid profile email users:read"

	testLargeNumericDate = 33305644800
	testSmallNumericDate = 1750497295

	testNonce  = "_Xm8-pWq29_LzR-k"
	testACR    = "urn:pxr:loa:1fa:any"
	testAZP    = "ELuXIbyntYQwIgUb8kFRmAotDzoGdPQZ"
	testAtHash = "U4osGiwnaxsnr4zCZBs2pA"
	testCHash  = "vF9dft4qmTGHjlOhiHcTqw"

	testClientID  = "test_client_id"
	testSessionID = "test_session_id"

	testName                = "John Michael Doe"
	testFamilyName          = "Doe"
	testGivenName           = "John"
	testMiddleName          = "Michael"
	testNickname            = "Johnny"
	testPreferredUsername   = "johndoe"
	testProfile             = "https://example.com/profile/johndoe"
	testPicture             = "https://example.com/profile/picture/johndoe"
	testWebsite             = "https://johndoe.dev"
	testGender              = "male"
	testBirthdate           = "1990-01-01"
	testZoneInfo            = "America/New_York"
	testLocale              = "en-US"
	testPhoneNumber         = "+1234567890"
	testPhoneNumberVerified = true
	testEmail               = "johndoe@example.com"
	testEmailVerified       = false
	testAddressFormatted    = "123 Main St\nMetropolis, NY 10001\nUSA"
	testStreetAddress       = "123 Main St"
	testLocality            = "Metropolis"
	testRegion              = "NY"
	testPostalCode          = "10001"
	testCountry             = "USA"

	testValidMinimalAccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJleHAiOjMzMzA1NjQ0ODAwLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20vIiwianRpIjoiYTBjNDNhYjMtMzA2Ny00MjExLTgwODYtZjZjN2YzMDA5YTgyIiwibmJmIjoxNzUwNDk3Mjk1LCJzdWIiOiIxIn0.C6WSp4EaynQBbdgLEJ1hpoHxwEoGAW8RYnCa1YpwNfE"
	testValidFullAccessToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiJ1cm46cHhyOmxvYToxZmE6YW55IiwiYW1yIjpbInB3ZCIsIm1mYSJdLCJhdF9oYXNoIjoiVTRvc0dpd25heHNucjR6Q1pCczJwQSIsImF1ZCI6InRlc3RBdWRpZW5jZSIsImF1dGhfdGltZSI6MTc1MDQ5NzI5NSwiYXpwIjoiRUx1WElieW50WVF3SWdVYjhrRlJtQW90RHpvR2RQUVoiLCJjX2hhc2giOiJ2RjlkZnQ0cW1UR0hqbE9oaUhjVHF3IiwiY2xpZW50X2lkIjoidGVzdF9jbGllbnRfaWQiLCJlbnRpdGxlbWVudHMiOlsidXNlcnM6cmVhZCIsInVzZXJzOndyaXRlIl0sImV4cCI6MzMzMDU2NDQ4MDAsImdyb3VwcyI6WyJhZG1pbiIsIm1vZGVyYXRvciJdLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20vIiwianRpIjoiYTBjNDNhYjMtMzA2Ny00MjExLTgwODYtZjZjN2YzMDA5YTgyIiwibmJmIjoxNzUwNDk3Mjk1LCJub25jZSI6Il9YbTgtcFdxMjlfTHpSLWsiLCJyb2xlcyI6WyJtZW1iZXIiLCJtYW5hZ2VyIl0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgdXNlcnM6cmVhZCIsInN1YiI6IjEifQ.C7-RKdKDhPPeVlz56fU2zZHma1mGi35iwPy2IrxWDBg"
	testInvalidAccessToken      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhd"
	testValidRefreshToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiJ1cm46cHhyOmxvYToxZmE6YW55IiwiYW1yIjpbInB3ZCIsIm1mYSJdLCJhdF9oYXNoIjoiVTRvc0dpd25heHNucjR6Q1pCczJwQSIsImF1ZCI6InRlc3RBdWRpZW5jZSIsImF1dGhfdGltZSI6MTc1MDQ5NzI5NSwiYXpwIjoiRUx1WElieW50WVF3SWdVYjhrRlJtQW90RHpvR2RQUVoiLCJjX2hhc2giOiJ2RjlkZnQ0cW1UR0hqbE9oaUhjVHF3IiwiY2xpZW50X2lkIjoidGVzdF9jbGllbnRfaWQiLCJlbnRpdGxlbWVudHMiOlsidXNlcnM6cmVhZCIsInVzZXJzOndyaXRlIl0sImV4cCI6MzMzMDU2NDQ4MDAsImdyb3VwcyI6WyJhZG1pbiIsIm1vZGVyYXRvciJdLCJpYXQiOjE3NTA0OTcyOTUsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20vIiwianRpIjoiZWY2YjVmNjAtNjcwMC00ZDA3LWEwYTctOWI1NjA0MzA1MzA0IiwibmJmIjoxNzUwNDk3Mjk1LCJub25jZSI6Il9YbTgtcFdxMjlfTHpSLWsiLCJyb2xlcyI6WyJtZW1iZXIiLCJtYW5hZ2VyIl0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgdXNlcnM6cmVhZCIsInNpZCI6InRlc3Rfc2Vzc2lvbl9pZCIsInN1YiI6IjEifQ.pVwJEoTpiPlxruv2QzSJ7RgIKOWXDuWn1yWtlf3Kro0"
	testValidIDToken            = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiJ1cm46cHhyOmxvYToxZmE6YW55IiwiYWRkcmVzcyI6eyJjb3VudHJ5IjoiVVNBIiwiZm9ybWF0dGVkIjoiMTIzIE1haW4gU3Rcbk1ldHJvcG9saXMsIE5ZIDEwMDAxXG5VU0EiLCJsb2NhbGl0eSI6Ik1ldHJvcG9saXMiLCJwb3N0YWxfY29kZSI6IjEwMDAxIiwicmVnaW9uIjoiTlkiLCJzdHJlZXRfYWRkcmVzcyI6IjEyMyBNYWluIFN0In0sImFtciI6WyJwd2QiLCJtZmEiXSwiYXRfaGFzaCI6IlU0b3NHaXduYXhzbnI0ekNaQnMycEEiLCJhdWQiOiJ0ZXN0QXVkaWVuY2UiLCJhdXRoX3RpbWUiOjE3NTA0OTcyOTUsImF6cCI6IkVMdVhJYnludFlRd0lnVWI4a0ZSbUFvdER6b0dkUFFaIiwiYmlydGhkYXRlIjoiMTk5MC0wMS0wMSIsImNfaGFzaCI6InZGOWRmdDRxbVRHSGpsT2hpSGNUcXciLCJlbWFpbCI6ImpvaG5kb2VAZXhhbXBsZS5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImV4cCI6MzMzMDU2NDQ4MDAsImZhbWlseV9uYW1lIjoiRG9lIiwiZ2VuZGVyIjoibWFsZSIsImdpdmVuX25hbWUiOiJKb2huIiwiaWF0IjoxNzUwNDk3Mjk1LCJpc3MiOiJodHRwczovL2V4YW1wbGUuY29tLyIsImp0aSI6IjNlZWRhNGJiLTk2ZDEtNDRiNi04ZWQyLWQ1OTc1N2M3ODZiZCIsImxvY2FsZSI6ImVuLVVTIiwibWlkZGxlX25hbWUiOiJNaWNoYWVsIiwibmFtZSI6IkpvaG4gTWljaGFlbCBEb2UiLCJuYmYiOjE3NTA0OTcyOTUsIm5pY2tuYW1lIjoiSm9obm55Iiwibm9uY2UiOiJfWG04LXBXcTI5X0x6Ui1rIiwicGhvbmVfbnVtYmVyIjoiKzEyMzQ1Njc4OTAiLCJwaG9uZV9udW1iZXJfdmVyaWZpZWQiOnRydWUsInBpY3R1cmUiOiJodHRwczovL2V4YW1wbGUuY29tL3Byb2ZpbGUvcGljdHVyZS9qb2huZG9lIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiam9obmRvZSIsInByb2ZpbGUiOiJodHRwczovL2V4YW1wbGUuY29tL3Byb2ZpbGUvam9obmRvZSIsInN1YiI6IjEiLCJ1cGRhdGVkX2F0IjoxNzUwNDk3Mjk1LCJ3ZWJzaXRlIjoiaHR0cHM6Ly9qb2huZG9lLmRldiIsInpvbmVpbmZvIjoiQW1lcmljYS9OZXdfWW9yayJ9.oeRZPFeTOe6dStBavZt1exy8_I_Wh-LpUB9Rlq87FTw"
)

var (
	testExpiry       = jwt.NewNumericDate(time.Unix(testLargeNumericDate, 0))
	testIssuedAt     = jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0))
	testNotBefore    = jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0))
	testAuthTime     = jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0))
	testAMR          = []string{"pwd", "mfa"}
	testRoles        = []string{"member", "manager"}
	testGroups       = []string{"admin", "moderator"}
	testEntitlements = []string{"users:read", "users:write"}
	testUpdatedAt    = jwt.NewNumericDate(time.Unix(testSmallNumericDate, 0))
	fixedDate        = time.Date(2026, 7, 15, 12, 30, 0, 0, time.UTC)
)

func NewTestMinimalAccessTokenClaims() claims.AccessTokenClaims {
	return claims.AccessTokenClaims{
		RegisteredClaims: newTestRegisteredClaims(testAccessTokenID),
	}
}

func NewTestFullAccessTokenClaims() claims.AccessTokenClaims {
	return claims.AccessTokenClaims{
		RegisteredClaims:     newTestRegisteredClaims(testAccessTokenID),
		AuthenticationClaims: newTestAuthenticationClaims(),
		AuthorizationClaims:  newTestAuthorizationClaims(),
		ClientClaims:         newTestClientClaims(),
	}
}

func NewTestMinimalRefreshTokenClaims() claims.RefreshTokenClaims {
	return claims.RefreshTokenClaims{
		RegisteredClaims: newTestRegisteredClaims(testRefreshTokenID),
	}
}

func NewTestFullRefreshTokenClaims() claims.RefreshTokenClaims {
	return claims.RefreshTokenClaims{
		RegisteredClaims:     newTestRegisteredClaims(testRefreshTokenID),
		AuthenticationClaims: newTestAuthenticationClaims(),
		AuthorizationClaims:  newTestAuthorizationClaims(),
		ClientClaims:         newTestClientClaims(),
		SessionClaims:        newTestSessionClaims(),
	}
}

func NewTestMinimalIDTokenClaims() claims.IDTokenClaims {
	return claims.IDTokenClaims{
		RegisteredClaims: newTestRegisteredClaims(testIDTokenID),
	}
}

func NewTestFullIDTokenClaims() claims.IDTokenClaims {
	return claims.IDTokenClaims{
		RegisteredClaims:        newTestRegisteredClaims(testIDTokenID),
		AuthenticationClaims:    newTestAuthenticationClaims(),
		OpenIDProfileClaims:     newTestOpenIDProfileClaims(),
		OpenIDPhoneNumberClaims: newTestOpenIDPhoneNumberClaims(),
		OpenIDEmailClaims:       newTestOpenIDEmailClaims(),
		OpenIDAddressClaims:     newTestOpenIDAddressClaims(),
	}
}

func newTestRegisteredClaims(id string) claims.RegisteredClaims {
	return claims.RegisteredClaims{
		ID:        id,
		Issuer:    testIssuer,
		Subject:   testSubject,
		Audience:  jwt.Audience{testAudience},
		Expiry:    testExpiry,
		IssuedAt:  testIssuedAt,
		NotBefore: testNotBefore,
	}
}

func newTestAuthenticationClaims() claims.AuthenticationClaims {
	return claims.AuthenticationClaims{
		AuthTime: testAuthTime,
		Nonce:    testNonce,
		ACR:      testACR,
		AMR:      testAMR,
		AZP:      testAZP,
		AtHash:   testAtHash,
		CHash:    testCHash,
	}
}

func newTestAuthorizationClaims() claims.AuthorizationClaims {
	return claims.AuthorizationClaims{
		Scope:        testScope,
		Roles:        testRoles,
		Groups:       testGroups,
		Entitlements: testEntitlements,
	}
}

func newTestClientClaims() claims.ClientClaims {
	return claims.ClientClaims{
		ClientID: testClientID,
	}
}

func newTestSessionClaims() claims.SessionClaims {
	return claims.SessionClaims{
		SessionID: testSessionID,
	}
}

func newTestOpenIDProfileClaims() claims.OpenIDProfileClaims {
	return claims.OpenIDProfileClaims{
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
		UpdatedAt:         testUpdatedAt,
	}
}

func newTestOpenIDPhoneNumberClaims() claims.OpenIDPhoneNumberClaims {
	verified := testPhoneNumberVerified

	return claims.OpenIDPhoneNumberClaims{
		PhoneNumber:         testPhoneNumber,
		PhoneNumberVerified: &verified,
	}
}

func newTestOpenIDEmailClaims() claims.OpenIDEmailClaims {
	verified := testEmailVerified

	return claims.OpenIDEmailClaims{
		Email:         testEmail,
		EmailVerified: &verified,
	}
}

func newTestOpenIDAddressClaims() claims.OpenIDAddressClaims {
	return claims.OpenIDAddressClaims{
		Address: &claims.AddressClaim{
			Formatted:     testAddressFormatted,
			StreetAddress: testStreetAddress,
			Locality:      testLocality,
			Region:        testRegion,
			PostalCode:    testPostalCode,
			Country:       testCountry,
		},
	}
}

func NewTestTokenWithRegisteredClaims() *jwt.JSONWebToken {
	token, err := jwt.ParseSigned(testValidMinimalAccessToken, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		panic(err)
	}

	return token
}

func NewTestAccessToken() *jwt.JSONWebToken {
	token, err := jwt.ParseSigned(testValidFullAccessToken, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		panic(err)
	}

	return token
}

func NewTestRefreshToken() *jwt.JSONWebToken {
	token, err := jwt.ParseSigned(testValidRefreshToken, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		panic(err)
	}

	return token
}

func NewTestIDToken() *jwt.JSONWebToken {
	token, err := jwt.ParseSigned(testValidIDToken, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		panic(err)
	}

	return token
}
