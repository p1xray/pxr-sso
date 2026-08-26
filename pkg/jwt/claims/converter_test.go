package claims

import (
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNumericDateToInt64(t *testing.T) {
	fixedTime := time.Date(2026, time.July, 21, 12, 0, 0, 0, time.UTC)
	expectedUnix := fixedTime.Unix()

	testCases := []struct {
		name     string
		input    *jwt.NumericDate
		expected int64
	}{
		{
			name:     "nil pointer returns zero",
			input:    nil,
			expected: 0,
		},
		{
			name:     "valid NumericDate returns correct int64 timestamp",
			input:    jwt.NewNumericDate(fixedTime),
			expected: expectedUnix,
		},
		{
			name:     "zero time NumericDate returns zero unix timestamp",
			input:    jwt.NewNumericDate(time.Time{}),
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := NumericDateToInt64(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
