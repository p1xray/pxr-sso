package validator

type Option func(*Validator)

func WithCustomClaims(f func() CustomClaims) Option {
	return func(v *Validator) {
		v.customClaims = f
	}
}
