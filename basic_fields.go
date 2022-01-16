package core

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.com/valyala/fastjson"
)

var (
	BasicFieldsInlineFields     = []string{"id", "ts_create", "ts_update", "removed", "info"}
	BasicFieldsInlineFieldsSize = len(BasicFieldsInlineFields)
)

// BasicFieldsInterface basic fields interface with dump function
type BasicFieldsInterface interface {
	Dump() []byte
	Name() string
}

// BasicFields basic fields
type BasicFields struct {
	BasicFieldsInterface

	ID       string    `json:"id"`
	TSCreate time.Time `json:"ts_create"`
	TSUpdate time.Time `json:"ts_update"`
	Removed  bool      `json:"removed"`
	Info     H         `json:"info"`
}

// NewBasicFields basicFields init
func NewBasicFields() BasicFields {
	now := Now()

	return BasicFields{
		ID:       xid.NewWithTime(now).String(),
		TSCreate: now,
		TSUpdate: now,
		Removed:  false,
		Info:     H{},
	}
}

// LoadBasicFields load with default fields
func LoadBasicFields(i interface{}, value *fastjson.Value) {
	x, ok := i.(*BasicFields)
	if !ok {
		return
	}

	x.ID = String(value.GetStringBytes("id"))
	x.TSCreate, _ = time.Parse(time.RFC3339Nano, String(value.GetStringBytes("ts_create")))
	x.TSUpdate, _ = time.Parse(time.RFC3339Nano, String(value.GetStringBytes("ts_update")))
	x.Removed = value.GetBool("removed")

	_ = jsoniter.Unmarshal(value.GetStringBytes("info"), &x.Info)
}

// PGXGet 获取单一对象
func (b *BasicFields) PGXGet(ctx context.Context, db *pgxpool.Pool, raw, order string) error {
	table_name := b.Name()
	fields_elm := reflect.ValueOf(b).Elem()
	fields_size := fields_elm.NumField()

	fields := make([]interface{}, 0, BasicFieldsInlineFieldsSize+fields_size)
	arguments := make([]string, 0, fields_size)
	values := make([]interface{}, 0, 1)

	if raw == "" {
		raw = "id=$1"
		values = append(values, b.ID)
	}

	if order == "" {
		order = "ts_create"
	}

	for i := 0; i < fields_size; i++ {
		valueField := fields_elm.Field(i)
		tag := fields_elm.Type().Field(i).Tag

		name := tag.Get("json")
		field := valueField.Addr().Interface()

		fields = append(fields, field)
		arguments = append(arguments, name)
	}

	return db.QueryRow(ctx, fmt.Sprintf(
		"select %s from %s where %s order by %s",
		strings.Join(arguments, ","), table_name, raw, order,
	), values...).Scan(fields...)
}

// PGXFilter 过滤
func (b *BasicFields) PGXFilter(ctx context.Context, db *pgxpool.Pool, raw, order string, scanWrapper func(pgx.Rows) error) error {
	table_name := b.Name()
	fields_elm := reflect.ValueOf(b).Elem()
	fields_size := fields_elm.NumField()

	if order == "" {
		order = "ts_create"
	}

	arguments := make([]string, 0, fields_size)

	for i := 0; i < fields_size; i++ {
		tag := fields_elm.Type().Field(i).Tag
		name := tag.Get("json")
		arguments = append(arguments, name)
	}

	rows, err := db.Query(ctx, fmt.Sprintf(
		"select %s from %s where %s order by %s",
		strings.Join(arguments, ","), table_name, raw, order,
	))
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := scanWrapper(rows); err != nil {
			return err
		}
	}

	return nil
}

// PGXCount 获取数量
func (b *BasicFields) PGXCount(ctx context.Context, db *pgxpool.Pool, raw string) (int, error) {
	table_name := b.Name()
	count := 0

	err := db.QueryRow(ctx, fmt.Sprintf("select count(*) from %s where %s", table_name, raw)).Scan(&count)
	return count, err
}

// PGXInsert 插入
func (b *BasicFields) PGXInsert(ctx context.Context, db *pgxpool.Pool) (pgx.Rows, error) {
	table_name := b.Name()
	fields_elm := reflect.ValueOf(b).Elem()
	fields_size := fields_elm.NumField()

	values := make([]interface{}, 0, BasicFieldsInlineFieldsSize+fields_size)

	values = append(values, b.ID)
	values = append(values, b.TSCreate)
	values = append(values, b.TSUpdate)
	values = append(values, b.Removed)
	values = append(values, b.Info)

	arguments := make([]string, 0, fields_size)
	arguments_q := make([]string, 0, fields_size)

	for i := 0; i < fields_size; i++ {
		valueField := fields_elm.Field(i)
		typeField := fields_elm.Type().Field(i)
		tag := typeField.Tag

		name := tag.Get("json")

		// ignore inline fields
		if sort.SearchStrings(BasicFieldsInlineFields, name) < BasicFieldsInlineFieldsSize {
			continue
		}

		values = append(values, valueField.Interface())
		arguments = append(arguments, name)
		arguments_q = append(arguments_q, fmt.Sprintf("$%d", len(values)))
	}

	return db.Query(ctx, fmt.Sprintf(
		"insert into %s (id,ts_create,ts_update,removed,info,%s) values($1,$2,$3,$4,$5,%s)",
		table_name, strings.Join(arguments, ","), strings.Join(arguments_q, ","),
	), values...)
}

// PGXUpdate 更新
func (b *BasicFields) PGXUpdate(ctx context.Context, db *pgxpool.Pool, data H) (pgx.Rows, error) {
	table_name := b.Name()
	fields_size := len(data)

	values := make([]interface{}, 0, fields_size)
	arguments := make([]string, 0, fields_size)

	// loop data
	for k, v := range data {
		if k == "id" {
			continue
		}

		values = append(values, v)
		arguments = append(arguments, fmt.Sprintf("%s=$%d", k, len(values)))
	}

	sql := fmt.Sprintf("update %s set %s where id=$1", table_name, strings.Join(arguments, ","))
	return db.Query(ctx, sql, values...)
}

// PGXRemove 删除
func (b *BasicFields) PGXRemove(ctx context.Context, db *pgxpool.Pool) error {
	table_name := b.Name()

	_, err := db.Exec(ctx, fmt.Sprintf("update %s set removed=true,ts_update=$2 where id=$1", table_name),
		b.ID, b.TSUpdate)

	return err
}

func init() {
	sort.Strings(BasicFieldsInlineFields)
}
