package core

import (
	"database/sql/driver"
	jsoniter "github.com/json-iterator/go"
)

// H alias map[string]any
type H map[string]any

// HNil empty struct
var HNil = H{}

func (h H) Fields() map[string]any {
	m := make(map[string]any, len(h))
	for k, v := range h {
		m[k] = v
	}
	return m
}

// Value implement driver.Valuer
func (h H) Value() (driver.Value, error) {
	return jsoniter.MarshalToString(h)
}

func (h H) LoadsString(src string) error {
	return jsoniter.UnmarshalFromString(src, &h)
}
func (h H) LoadsBytes(src []byte) error {
	return jsoniter.Unmarshal(src, &h)
}
