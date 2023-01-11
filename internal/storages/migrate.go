package storages

import (
	"context"

	"github.com/jackc/tern/migrate"
	"github.com/pkg/errors"

	"GameJamPlatform/internal/log"
)

func (st *storage) Migrate(ctx context.Context) error {
	if !st.cfg.MigrationEnabled {
		return nil
	}

	migrator, err := migrate.NewMigrator(ctx, st.db, "version")
	if err != nil {
		return errors.Wrap(err, "Migrator error")
	}

	err = migrator.LoadMigrations(st.cfg.MigrationPath)
	if err != nil {
		return errors.Wrap(err, "Migrator error")
	}

	ver, err := migrator.GetCurrentVersion(ctx)

	if ver == 0 {
		log.Debug("Migrate to last version")

		err = migrator.Migrate(ctx)
		if err != nil {
			return errors.Wrap(err, "Migrator error")
		}

		log.Debug("Migrate to last version success")
		return nil
	}

	if st.cfg.MigrationVersion > 0 && st.cfg.MigrationVersion != ver {
		log.Debug("Migrate to version %d", st.cfg.MigrationVersion)

		err = migrator.MigrateTo(ctx, st.cfg.MigrationVersion)
		if err != nil {
			return errors.Wrap(err, "Migrator error")
		}

		log.Debug("Migrate to version %d success", st.cfg.MigrationVersion)
		return nil
	}

	return nil
}
