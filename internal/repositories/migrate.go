package repositories

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/meehighlov/workout/internal/config"
	"github.com/meehighlov/workout/internal/migrations"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func migrateWithGoose(db *gorm.DB, cfg *config.Config, logger *slog.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get *sql.DB from GORM", "error", err)
		return err
	}

	goose.WithLogger(&gooseLogger{logger: logger})

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Failed to set Goose dialect", "error", err)
		return err
	}

	tempDir, err := os.MkdirTemp("", cfg.MigrationsDir)
	if err != nil {
		logger.Error("Failed to create temp directory for migrations", "error", err)
		return err
	}
	defer os.RemoveAll(tempDir)

	logger.Info("Extracting embedded migrations to temp directory", "dir", tempDir)
	if err := extractEmbeddedMigrations(tempDir, logger); err != nil {
		return err
	}

	logger.Info("Running SQL migrations with Goose from embedded files")
	if err := goose.Up(sqlDB, tempDir); err != nil {
		logger.Error("Failed to run Goose migrations", "error", err)
		return err
	}

	logger.Info("Embedded SQL migrations completed successfully")
	return nil
}

func extractEmbeddedMigrations(tempDir string, logger *slog.Logger) error {
	entries, err := fs.ReadDir(migrations.SQLMigrations, ".")
	if err != nil {
		logger.Error("Failed to read embedded migrations directory", "error", err)
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if filepath.Ext(fileName) != ".sql" {
			continue
		}

		data, err := fs.ReadFile(migrations.SQLMigrations, fileName)
		if err != nil {
			logger.Error("Failed to read embedded migration file", "file", fileName, "error", err)
			return err
		}

		destPath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			logger.Error("Failed to write migration file to temp directory", "file", fileName, "error", err)
			return err
		}

		logger.Debug("Extracted migration file", "file", fileName)
	}

	return nil
}

func RunMigrations(ctx context.Context, cfg *config.Config, logger *slog.Logger, db *gorm.DB) error {
	logger.Info("Starting database migrations")

	db = db.WithContext(ctx)

	if cfg.RunMigrations {
		if err := migrateWithGoose(db, cfg, logger); err != nil {
			return err
		}
	}

	logger.Info("All migrations completed successfully")
	return nil
}

type gooseLogger struct {
	logger *slog.Logger
}

func (l *gooseLogger) Fatalf(format string, v ...interface{}) {
	msg := format
	if len(v) > 0 {
		msg = sprintf(format, v...)
	}
	l.logger.Error("Goose fatal error", "message", msg)
	panic(msg)
}

func (l *gooseLogger) Printf(format string, v ...interface{}) {
	msg := format
	if len(v) > 0 {
		msg = sprintf(format, v...)
	}
	l.logger.Info("Goose", "message", msg)
}

func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
