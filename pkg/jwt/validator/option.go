package validator

import "time"

type Option func(*Validator)

func WithCustomClaims(f func() CustomClaims) Option {
	return func(v *Validator) {
		v.customClaims = f
	}
}

func WithAllowedClockSkew(d time.Duration) Option {
	return func(v *Validator) {
		v.allowedClockSkew = d
	}
}
