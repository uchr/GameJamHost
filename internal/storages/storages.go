package storages

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type storage struct {
	db  *pgx.Conn
	cfg *Config
}

var _ Repo = (*storage)(nil)

func NewStorage(ctx context.Context, cfg *Config) (*storage, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	db, err := pgx.Connect(ctxWithTimeout, cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	st := &storage{db: db, cfg: cfg}

	if err := st.Migrate(ctx); err != nil {
		return nil, err
	}

	return st, nil
}
