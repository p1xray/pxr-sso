package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres provides access to PostgreSQL connections.
type Postgres struct {
	maxPoolSize        int
	connAttempts       int
	connAttemptTimeout time.Duration

	Pool *pgxpool.Pool
}

// New creates a new instance of the PostgreSQL connections.
func New(url string, setters ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:        _defaultMaxPoolSize,
		connAttempts:       _defaultConnAttempts,
		connAttemptTimeout: _defaultConnTimeout,
	}

	for _, setter := range setters {
		setter(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		time.Sleep(pg.connAttemptTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	return pg, nil
}

// Close closes connections to PostgreSQL.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
