package core

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PGXGet 获取单一对象
func PGXGet[T BasicFieldsInterface](ctx context.Context, db *pgxpool.Pool,
	item T, raw, order string, key any,
) error {
	table := item.TableName()
	fields, arguments := basicFieldsGenFieldsAndArguments(item)
	query := buildSelectQuery(table, raw, order, arguments...)

	err := db.QueryRow(ctx, query, key).Scan(fields...)
	if err == pgx.ErrNoRows {
		return ErrObjectNotFound
	}
	return err
}

// PGXFilter 过滤
func PGXFilter[T BasicFieldsInterface](ctx context.Context, db *pgxpool.Pool,
	item T, raw, order string, offset, limit int64, scanWrapper func(pgx.Rows) error, params ...any,
) error {
	_, arguments := basicFieldsGenFieldsAndArguments(item)
	table := item.TableName()
	query := buildSelectQuery(table, raw, order, arguments...)

	builder := strings.Builder{}
	builder.WriteString(query)
	builder.WriteString(" LIMIT ")
	builder.WriteString(strconv.FormatInt(limit, 10))
	builder.WriteString(" OFFSET ")
	builder.WriteString(strconv.FormatInt(offset, 10))

	rows, err := db.Query(ctx, builder.String(), params...)
	if err != nil {
		E("query error", err, H{"sql": query})
		return err
	}

	defer rows.Close()

	for rows.Next() {
		if err := scanWrapper(rows); err != nil {
			return err
		}
	}

	return nil
}

// PGXInsert 插入
func PGXInsert[T BasicFieldsInterface](ctx context.Context, db *pgxpool.Pool, item T) error {
	values, arguments := basicFieldsGenFieldsAndArguments(item)
	table := item.TableName()

	argumentsQ := make([]string, len(arguments))
	for i := range arguments {
		argumentsQ[i] = "$" + strconv.Itoa(i+1)
	}

	query := strings.Builder{}
	query.WriteString("INSERT INTO ")
	query.WriteString(table)
	query.WriteString(" (")
	query.WriteString(strings.Join(arguments, ","))
	query.WriteString(") VALUES (")
	query.WriteString(strings.Join(argumentsQ, ","))
	query.WriteString(")")

	_, err := db.Exec(ctx, query.String(), values...)
	return err
}

// PGXUpdate 更新
func PGXUpdate[T BasicFieldsInterface](ctx context.Context, db *pgxpool.Pool, item T, key any, data H) error {
	size := len(data)
	values := make([]any, 0, size)
	arguments := make([]string, 0, size)

	values = append(values, key)

	builder := strings.Builder{}

	// loop data
	for k, v := range data {
		if k == "id" {
			continue
		}

		values = append(values, v)

		builder.Reset()
		builder.WriteString(k)
		builder.WriteString("=$")
		builder.WriteString(strconv.Itoa(len(values)))

		arguments = append(arguments, builder.String())
	}

	builder.Reset()
	builder.WriteString("UPDATE ")
	builder.WriteString(item.TableName())
	builder.WriteString(" SET ")
	builder.WriteString(strings.Join(arguments, ", "))
	builder.WriteString(" WHERE id=$1")

	_, err := db.Exec(ctx, builder.String(), values...)
	return err
}

// PGXCount 计数
func PGXCount(ctx context.Context, db *pgxpool.Pool, item BasicFieldsInterface, raw string, params ...any) (int64, error) {
	var count int64

	if raw != "" {
		raw = "where " + raw
	}

	query := strings.Builder{}
	query.WriteString("SELECT COUNT(*) FROM ")
	query.WriteString(item.TableName())
	query.WriteString(" ")
	query.WriteString(raw)

	err := db.QueryRow(ctx, query.String(), params...).Scan(&count)
	return count, err
}

// PGXRemove 删除
func PGXRemove(ctx context.Context, db *pgxpool.Pool, item BasicFieldsInterface, key any) error {
	query := strings.Builder{}
	query.WriteString("UPDATE ")
	query.WriteString(item.TableName())
	query.WriteString(" SET ")
	query.WriteString("removed=true, ts_update=$2")
	query.WriteString(" WHERE id=$1 and not removed")

	_, err := db.Exec(ctx, query.String(), key, Now())
	return err
}

// PGXUpdateHelper 快速更新的辅佐
func PGXUpdateHelper(ctx context.Context, name string, id any, data H, db *pgxpool.Pool) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	for key, value := range data {
		// query := fmt.Sprintf("update %s set %s=$2 where id=$1", name, key)
		builder.Reset()
		builder.WriteString("update ")
		builder.WriteString(name)
		builder.WriteString(" set ")
		builder.WriteString(key)
		builder.WriteString("=$2 where id=$1")

		if _, err := tx.Exec(ctx, builder.String(), id, value); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}
