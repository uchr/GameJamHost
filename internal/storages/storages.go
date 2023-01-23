package storages

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type storage struct {
	db  *pgxpool.Pool
	cfg *Config
}

var _ Repo = (*storage)(nil)

func NewStorage(ctx context.Context, cfg *Config) (*storage, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	db, err := pgxpool.Connect(ctxWithTimeout, cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	st := &storage{
		db:  db,
		cfg: cfg,
	}

	if err := st.migrate(ctx, cfg.DatabasePath, cfg.ConnectTimeout); err != nil {
		return nil, err
	}

	return st, nil
}
