package core

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PGXUpdateHelper simple update query
func PGXUpdateHelper(ctx context.Context, name string, data H, db *pgxpool.Pool) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	for key, name := range data {
		query := fmt.Sprintf("update %s set %s=$2 where id=$1", name, key)

		if _, err := tx.Exec(ctx, query, name); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

	}

	return tx.Commit(ctx)
}
