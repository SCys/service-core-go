package core

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
)

// DBGet get one row from database
func DBGet[T BasicFieldsInterface](db *sql.DB, item T, raw, order string, key any) error {
	table := item.TableName()
	fields, arguments := basicFieldsGenFieldsAndArguments(item)

	query := buildSelectQuery(table, raw, order, arguments...)

	err := db.QueryRow(query, key).Scan(fields...)
	if err == pgx.ErrNoRows {
		return ErrObjectNotFound
	}

	return err
}

// DBFilter filter rows from database
func DBFilter[T BasicFieldsInterface](db *sql.DB,
	item T,
	raw, order string, offset, limit int64, scanWrapper func(*sql.Rows) error, params ...any,
) error {
	_, arguments := basicFieldsGenFieldsAndArguments(item)
	table := item.TableName()
	query := buildSelectQuery(table, raw, order, arguments...)

	rows, err := db.Query(query+fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset), params...)
	if err != nil {
		E("query error", err, H{"sql": query})
		return err
	}

	for rows.Next() {
		if err := scanWrapper(rows); err != nil {
			return err
		}
	}

	return nil
}

// DBInsert insert one row to database
func DBInsert[T BasicFieldsInterface](db *sql.DB, item T) error {
	values, arguments := basicFieldsGenFieldsAndArguments(item)
	table := item.TableName()

	argumentsQ := make([]string, len(arguments))
	for i := range arguments {
		argumentsQ[i] = fmt.Sprintf("$%d", i+1)
	}

	query := strings.Builder{}
	query.WriteString("INSERT INTO ")
	query.WriteString(table)
	query.WriteString(" (")
	query.WriteString(strings.Join(arguments, ","))
	query.WriteString(") VALUES (")
	query.WriteString(strings.Join(argumentsQ, ","))
	query.WriteString(")")

	_, err := db.Exec(query.String(), values...)
	return err
}

// DBUpdate update row in database
func DBUpdate[T BasicFieldsInterface](db *sql.DB, item T, key any, data H) error {
	size := len(data)
	values := make([]any, 0, size)
	arguments := make([]string, 0, size)

	values = append(values, key)

	// loop data
	for k, v := range data {
		if k == "id" {
			continue
		}

		values = append(values, v)
		arguments = append(arguments, fmt.Sprintf("%s=$%d", k, len(values)))
	}

	query := strings.Builder{}
	query.WriteString("UPDATE ")
	query.WriteString(item.TableName())
	query.WriteString(" SET ")
	query.WriteString(strings.Join(arguments, ", "))
	query.WriteString(" WHERE id=$1")

	_, err := db.Exec(query.String(), values...)
	return err
}

// DBCount count rows in database
func DBCount[T BasicFieldsInterface](db *sql.DB, item T, raw string) (int, error) {
	count := 0

	if raw != "" {
		raw = "where " + raw
	}

	query := strings.Builder{}
	query.WriteString("SELECT COUNT(*) FROM ")
	query.WriteString(item.TableName())
	query.WriteString(" ")
	query.WriteString(raw)

	err := db.QueryRow(query.String()).Scan(&count)
	return count, err
}

// DBRemove remove row from database
func DBRemove[T BasicFieldsInterface](db *sql.DB, item T, key any) error {
	query := strings.Builder{}
	query.WriteString("UPDATE ")
	query.WriteString(item.TableName())
	query.WriteString(" SET ")
	query.WriteString("removed=true, ts_update=$2")
	query.WriteString(" WHERE id=$1 and not removed")

	_, err := db.Exec(query.String(), key, Now())
	return err
}
