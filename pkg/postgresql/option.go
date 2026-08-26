package postgresql

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}

func ConnAttemptTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		p.connAttemptTimeout = timeout
	}
}
