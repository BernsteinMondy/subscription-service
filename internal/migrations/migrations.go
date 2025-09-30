package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func Run(ctx context.Context, db *sql.DB, dir string) error {
	err := goose.SetDialect(string(goose.DialectPostgres))
	if err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	err = goose.UpContext(ctx, db, dir)
	if err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}

	return nil
}
