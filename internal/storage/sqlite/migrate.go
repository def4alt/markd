package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set goose sqlite dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("run sqlite migrations: %w", err)
	}

	return nil
}
