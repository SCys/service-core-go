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
	isPointer := targetValue.Kind() == reflect.Pointer

	if isPointer {
		targetValue = targetValue.Elem()
		targetType = targetType.Elem()
	}

	fieldsSize := targetType.NumField()

	fields := make([]any, 0, basicFieldsInlineFieldsSize+fieldsSize)
	arguments := make([]string, 0, basicFieldsInlineFieldsSize+fieldsSize)

	for i := 0; i < fieldsSize; i++ {
		v := targetValue.Field(i)
		t := targetType.Field(i)

		// 解码基础 BasicFields 对象
		if v.Type().Kind() == reflect.Struct && v.Type().Name() == "BasicFields" {
			if isPointer {
				fields = append(fields, targetValue.FieldByName("ID").Addr().Interface())
				fields = append(fields, targetValue.FieldByName("TSCreate").Addr().Interface())
				fields = append(fields, targetValue.FieldByName("TSUpdate").Addr().Interface())
				fields = append(fields, targetValue.FieldByName("Removed").Addr().Interface())
				fields = append(fields, targetValue.FieldByName("Info").Addr().Interface())
			} else {
				fields = append(fields, targetValue.FieldByName("ID").Interface())
				fields = append(fields, targetValue.FieldByName("TSCreate").Interface())
				fields = append(fields, targetValue.FieldByName("TSUpdate").Interface())
				fields = append(fields, targetValue.FieldByName("Removed").Interface())
				fields = append(fields, targetValue.FieldByName("Info").Interface())
			}

			arguments = append(arguments, "id")
			arguments = append(arguments, "ts_create")
			arguments = append(arguments, "ts_update")
			arguments = append(arguments, "removed")
			arguments = append(arguments, "info")

			continue
		}

		// only support json tag
		tag := t.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		name := strings.Split(tag, ",")[0]

		// 修复指针类型字段处理
		if isPointer {
			// 对于指针类型，我们需要获取字段的地址
			fieldAddr := targetValue.Field(i).Addr()
			fields = append(fields, fieldAddr.Interface())
		} else {
			// 对于值类型，直接获取字段值
			fields = append(fields, targetValue.Field(i).Interface())
		}
		arguments = append(arguments, name)
	}

	return fields, arguments
}
