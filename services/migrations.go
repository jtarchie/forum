package services

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jtarchie/forum/db"
	"go.uber.org/zap"
)

//go:embed assets/migration.sql
var migrationSQL string

//go:embed assets/migrations/*.sql
var migrations embed.FS

func Migration(client db.Client, logger *zap.Logger) error {
	err := client.Execute(migrationSQL)
	if err != nil {
		return fmt.Errorf("could not execute migrations: %w", err)
	}

	matches, err := fs.Glob(migrations, "assets/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("could not find matches for migrations: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("no migrations were found")
	}

	for _, migration := range matches {
		version := strings.Split(filepath.Base(migration), ".")[0]

		logger.Info("migration.found", zap.String("version", version))

		rows, err := client.Query("SELECT 1 FROM migrations WHERE version = ?", version)
		if err != nil {
			return fmt.Errorf("could not read migration %q: %w", version, err)
		}

		count := 0
		for rows.Next() {
			count++
		}

		if count == 0 {
			logger.Info("migration.missing", zap.String("version", version))

			content, err := migrations.ReadFile(migration)
			if err != nil {
				return fmt.Errorf("could not read migration file (%q): %w", migration, err)
			}

			err = client.Execute(string(content))
			if err != nil {
				return fmt.Errorf("could not run migrations (%q): %w", migration, err)
			}

			err = client.Execute("INSERT INTO migrations (version) VALUES (?);", version)
			if err != nil {
				return fmt.Errorf("could not set migration (%q): %w", migration, err)
			}
		} else {
			logger.Info("migration.exists", zap.String("version", version))
		}
	}

	return nil
}
