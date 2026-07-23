package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testIssuer  = "https://example.com/"
	testSubject = "1"
	leeway      = 1 * time.Minute
)

var (
	testAudience = jwt.Audience{"test_audience", "other_audience"}
	testTTL      = 1 * time.Hour
)

func TestNewRegisteredClaims(t *testing.T) {
	testCases := []struct {
		name     string
		issuer   string
		subject  string
		audience jwt.Audience
		ttl      time.Duration
	}{
		{
			name:     "successful creation of new registered claims",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      1 * time.Hour,
		},
		{
			name:     "successful creation of new registered claims with zero TTL",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      0,
		},
		{
			name:     "successful creation of new registered claims with negative TTL",
			issuer:   testIssuer,
			subject:  testSubject,
			audience: testAudience,
			ttl:      -5 * time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registeredClaims := NewRegisteredClaims(tc.issuer, tc.subject, tc.audience, tc.ttl)

			checkRegisteredClaims(t, registeredClaims, tc.issuer, tc.subject, tc.audience, leeway)
		})
	}
}

func checkRegisteredClaims(
	t *testing.T,
	registeredClaims RegisteredClaims,
	expectedIssuer string,
	expectedSubject string,
	expectedAudience jwt.Audience,
	leeway time.Duration,
) {
	assert.NotEmpty(t, registeredClaims.ID)
	_, err := uuid.Parse(registeredClaims.ID)
	assert.NoError(t, err)

	assert.Equal(t, registeredClaims.Issuer, expectedIssuer)
	assert.Equal(t, registeredClaims.Subject, expectedSubject)
	assert.Equal(t, registeredClaims.Audience, expectedAudience)

	assert.True(t, time.Now().Add(-leeway).Before(registeredClaims.Expiry.Time()))
	assert.True(t, time.Now().Add(leeway).After(registeredClaims.NotBefore.Time()))
	assert.True(t, time.Now().Add(leeway).After(registeredClaims.IssuedAt.Time()))
}
