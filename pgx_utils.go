package core

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PGXUpdateHelper 快速更新的辅佐
func PGXUpdateHelper(ctx context.Context, name string, id interface{}, data H, db *pgxpool.Pool) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	for key, value := range data {
		query := fmt.Sprintf("update %s set %s=$2 where id=$1", name, key)
		if _, err := tx.Exec(ctx, query, id, value); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}
