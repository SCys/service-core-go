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
}

// BasicFields basic fields
type BasicFields struct {
	ID       string    `json:"id"`
	TSCreate time.Time `json:"ts_create"`
	TSUpdate time.Time `json:"ts_update"`
	Removed  bool      `json:"removed"`
	Info     H         `json:"info"`
}

// Dump dump to bytes
func (b *BasicFields) Dump() []byte {
	raw, err := jsoniter.Marshal(b)
	if err != nil {
		E("json unmarshal error", err)
	}
	return raw
}

// PGXGet 获取单一对象
func PGXGet(ctx context.Context, b interface{}, table string, db *pgxpool.Pool, raw, order string, key interface{}) error {
	fieldsElm := reflect.ValueOf(b).Elem()
	fieldsSize := fieldsElm.NumField()

	fields := make([]interface{}, 0, BasicFieldsInlineFieldsSize+fieldsSize)
	arguments := make([]string, 0, BasicFieldsInlineFieldsSize+fieldsSize)

	if key == "" {
		key = b.(*BasicFields).ID
	}

	if raw == "" {
		raw = "id=$1"
	}

	if order == "" {
		order = "ts_create"
	}

	for i := 0; i < fieldsSize; i++ {
		v := fieldsElm.Field(i)
		t := fieldsElm.Type().Field(i)

		if v.Type().Kind() == reflect.Struct && t.Anonymous && v.Type().Name() == "BasicFields" {
			fieldEmbedElm := reflect.ValueOf(v.Addr().Interface()).Elem()
			for j := 0; j < fieldEmbedElm.NumField(); j++ {
				arguments = append(arguments, fieldEmbedElm.Type().Field(j).Tag.Get("json"))
				fields = append(fields, fieldEmbedElm.Field(j).Addr().Interface())
			}

			continue
		}

		// only support jsoniterable
		name := t.Tag.Get("json")
		if name == "" {
			continue
		}

		fields = append(fields, v.Addr().Interface())
		arguments = append(arguments, name)
	}

	sql := fmt.Sprintf("select %s from %s where %s order by %s", strings.Join(arguments, ","), table, raw, order)

	// D("sql: %s", sql)
	// D("fields: %v", fields)

	err := db.QueryRow(ctx, sql, key).Scan(fields...)
	if err == pgx.ErrNoRows {
		return ErrObjectNotFound
	}

	return err
}

// PGXFilter 过滤
func PGXFilter(
	ctx context.Context, b interface{}, tableName string, db *pgxpool.Pool,
	raw, order string, offset, limit int64, scanWrapper func(pgx.Rows) error, params ...interface{},
) error {
	fieldsElm := reflect.ValueOf(b).Elem()
	fieldsSize := fieldsElm.NumField()

	if order == "" {
		order = "ts_create"
	}

	arguments := make([]string, 0, fieldsSize)

	for i := 0; i < fieldsSize; i++ {
		v := fieldsElm.Field(i)
		t := fieldsElm.Type().Field(i)

		if v.Type().Kind() == reflect.Struct && t.Anonymous && v.Type().Name() == "BasicFields" {
			fieldEmbedElm := reflect.ValueOf(v.Addr().Interface()).Elem()
			for j := 0; j < fieldEmbedElm.NumField(); j++ {
				arguments = append(arguments, fieldEmbedElm.Type().Field(j).Tag.Get("json"))
			}

			continue
		}

		arguments = append(arguments, t.Tag.Get("json"))
	}

	sql := fmt.Sprintf(
		"select %s from %s where %s order by %s offset %d limit %d",
		strings.Join(arguments, ","), tableName, raw, order, offset, limit,
	)
	rows, err := db.Query(ctx, sql, params...)
	if err != nil {
		E("query error", err, H{"sql": sql})
		return err
	}

	for rows.Next() {
		if err := scanWrapper(rows); err != nil {
			return err
		}
	}

	return nil
}

// PGXInsert 插入
func PGXInsert(ctx context.Context, b interface{}, tableName string, db *pgxpool.Pool) error {
	fieldsElm := reflect.ValueOf(b).Elem()
	fieldsSize := fieldsElm.NumField()

	values := make([]interface{}, 0, BasicFieldsInlineFieldsSize+fieldsSize)

	arguments := make([]string, 0, BasicFieldsInlineFieldsSize+fieldsSize)
	argumentsQ := make([]string, 0, BasicFieldsInlineFieldsSize+fieldsSize)

	for i := 0; i < fieldsSize; i++ {
		v := fieldsElm.Field(i)
		t := fieldsElm.Type().Field(i)

		if v.Type().Kind() == reflect.Struct && t.Anonymous && v.Type().Name() == "BasicFields" {
			fieldEmbedElm := reflect.ValueOf(v.Addr().Interface()).Elem()
			for j := 0; j < fieldEmbedElm.NumField(); j++ {
				values = append(values, fieldEmbedElm.Field(j).Addr().Interface())
				arguments = append(arguments, fieldEmbedElm.Type().Field(j).Tag.Get("json"))
				argumentsQ = append(argumentsQ, fmt.Sprintf("$%d", len(values)))
			}

			continue
		}

		name := t.Tag.Get("json")

		values = append(values, v.Interface())
		arguments = append(arguments, name)
		argumentsQ = append(argumentsQ, fmt.Sprintf("$%d", len(values)))
	}

	_, err := db.Exec(ctx, fmt.Sprintf(
		"insert into %s (%s) values(%s)",
		tableName, strings.Join(arguments, ","), strings.Join(argumentsQ, ","),
	), values...)
	return err
}

// PGXUpdate 更新
func PGXUpdate(ctx context.Context, b interface{}, tableName string, db *pgxpool.Pool, data H) (pgx.Rows, error) {
	size := len(data)
	values := make([]interface{}, 0, size)
	arguments := make([]string, 0, size)

	// loop data
	for k, v := range data {
		if k == "id" {
			continue
		}

		values = append(values, v)
		arguments = append(arguments, fmt.Sprintf("%s=$%d", k, len(values)))
	}

	sql := fmt.Sprintf("update %s set %s where id=$1", tableName, strings.Join(arguments, ","))
	return db.Query(ctx, sql, values...)
}

// PGXCount 获取数量
func (b *BasicFields) PGXCount(ctx context.Context, tableName string, db *pgxpool.Pool, raw string) (int, error) {
	count := 0

	if raw != "" {
		raw = "where " + raw
	}

	err := db.QueryRow(ctx, fmt.Sprintf("select count(*) from %s %s", tableName, raw)).Scan(&count)
	return count, err
}

// PGXRemove 删除
func (b *BasicFields) PGXRemove(ctx context.Context, tableName string, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, fmt.Sprintf("update %s set removed=true,ts_update=$2 where id=$1", tableName),
		b.ID, b.TSUpdate)

	return err
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

func init() {
	sort.Strings(BasicFieldsInlineFields)
}
