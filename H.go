package core

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
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
	return sonic.MarshalString(h)
}

func (h H) LoadsString(src string) error {
	return sonic.UnmarshalString(src, &h)
}
func (h H) LoadsBytes(src []byte) error {
	return sonic.Unmarshal(src, &h)
}
