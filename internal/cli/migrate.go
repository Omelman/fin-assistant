package cli

import (
	"github.com/fin-assistant/internal/db-assets"
	"github.com/fin-assistant/internal/config"
	migrate "github.com/rubenv/sql-migrate"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var migrations = &migrate.PackrMigrationSource{
	Box: db_assets.Migrations,
}

func MigrateUp(cfg config.Config) error {
	applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	cfg.Log().WithField("applied", applied).Info("migrations applied")
	return nil
}

func MigrateDown(cfg config.Config) error {
	applied, err := migrate.Exec(cfg.DB().RawDB(), "implementation", migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	cfg.Log().WithField("applied", applied).Info("migrations applied")
	return nil
}
