package storages

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"

	"GameJamPlatform/internal/log"
)

func (st *storage) migrate(ctx context.Context, databasePath string, connectionTimeout time.Duration) error {
	if !st.cfg.MigrationEnabled {
		return nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	conn, err := pgx.Connect(ctxWithTimeout, databasePath)
	if err != nil {
		return fmt.Errorf("migrator error: %w", err)
	}

	migrator, err := migrate.NewMigrator(ctx, conn, "version")
	if err != nil {
		return fmt.Errorf("migrator error: %w", err)
	}

	err = migrator.LoadMigrations(st.cfg.MigrationPath)
	if err != nil {
		return fmt.Errorf("migrator error: %w", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)

	if ver == 0 {
		log.Debug("Migrate to last version")

		err = migrator.Migrate(ctx)
		if err != nil {
			return fmt.Errorf("migrator error: %w", err)
		}

		log.Debug("Migrate to last version success")
		return nil
	}

	if st.cfg.MigrationVersion > 0 && st.cfg.MigrationVersion != ver {
		log.Debug("Migrate to version %d", st.cfg.MigrationVersion)

		err = migrator.MigrateTo(ctx, st.cfg.MigrationVersion)
		if err != nil {
			return fmt.Errorf("migrator error: %w", err)
		}

		log.Debug("Migrate to version %d success", st.cfg.MigrationVersion)
		return nil
	}

	return nil
}
