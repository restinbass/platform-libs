package db_test_tools

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDirName   = "migrations"
	postgresDialectName = "postgres"
)

var PgxPool *pgxpool.Pool

// UpMigrations -
func UpMigrations(ctx context.Context) error {
	migrationsDirPath, err := findMigrationsDirPath()
	if err != nil {
		return err
	}

	if err := goose.SetDialect(postgresDialectName); err != nil {
		return err
	}

	goose.SetVerbose(false)
	PgxPool, err = pgxpool.New(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		return err
	}

	if err := PgxPool.Ping(ctx); err != nil {
		return err
	}

	return goose.Up(stdlib.OpenDBFromPool(PgxPool), migrationsDirPath)
}

// DownMigrations -
func DownMigrations() error {
	if PgxPool == nil {
		return nil
	}

	migrationsDirPath, err := findMigrationsDirPath()
	if err != nil {
		return err
	}

	return goose.Reset(stdlib.OpenDB(*PgxPool.Config().ConnConfig), migrationsDirPath)
}

func findMigrationsDirPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		currentDirectory := filepath.Join(dir, migrationsDirName)
		fileInfo, err := os.Stat(currentDirectory)
		if err != nil {
			if os.IsNotExist(err) {
				dir = filepath.Dir(dir)
				continue
			}
			return "", err
		}

		if !fileInfo.IsDir() {
			continue
		}

		return currentDirectory, nil
	}
}
