package core

import (
	"reflect"
	"strings"
)

func buildSelectQuery(table, raw, order string, arguments ...string) string {
	query := strings.Builder{}
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(arguments, ","))
	query.WriteString(" FROM ")
	query.WriteString(table)

	if raw != "" {
		query.WriteString(" WHERE ")
		query.WriteString(raw)
	}

	if order != "" {
		query.WriteString(" ORDER BY ")
		query.WriteString(order)
	}

	return query.String()
}

func basicFieldsGenFieldsAndArguments[T BasicFieldsInterface](target T) ([]any, []string) {
	targetValue := reflect.ValueOf(target)
	targetType := reflect.TypeOf(target)

	if targetValue.Kind() == reflect.Pointer {
		targetValue = targetValue.Elem()
		targetType = targetType.Elem()
	}

	fieldsSize := targetType.NumField()

	values := make([]any, 0, basicFieldsInlineFieldsSize+fieldsSize)
	arguments := make([]string, 0, basicFieldsInlineFieldsSize+fieldsSize)

	for i := 0; i < fieldsSize; i++ {
		v := targetValue.Field(i)
		t := targetType.Field(i)

		// 解码基础 BasicFields 对象
		if v.Type().Kind() == reflect.Struct && v.Type().Name() == "BasicFields" {
			values = append(values, targetValue.FieldByName("ID").Interface())
			values = append(values, targetValue.FieldByName("TSCreate").Interface())
			values = append(values, targetValue.FieldByName("TSUpdate").Interface())
			values = append(values, targetValue.FieldByName("Removed").Interface())
			values = append(values, targetValue.FieldByName("Info").Interface())

			arguments = append(arguments, "id")
			arguments = append(arguments, "ts_create")
			arguments = append(arguments, "ts_update")
			arguments = append(arguments, "removed")
			arguments = append(arguments, "info")

			continue
		}

		// only support json tag
		tag := t.Tag.Get("json")
		if tag == "" {
			continue
		}

		name := strings.Split(tag, ",")[0]

		values = append(values, v.Interface())
		arguments = append(arguments, name)
	}

	return values, arguments
}
