package core

import (
	"sort"
	"time"

	"github.com/rs/xid"

	sonic "github.com/bytedance/sonic"
)

// 基础字段
var (
	basicFieldsInlineFields     = []string{"id", "ts_create", "ts_update", "removed", "info"}
	basicFieldsInlineFieldsSize = len(basicFieldsInlineFields)
)

// BasicFieldsInterface basic fields interface with dump function
type BasicFieldsInterface interface {
	TableName() string
}

// BasicFields basic fields
type BasicFields struct {
	ID       string    `json:"id"`
	TSCreate time.Time `json:"ts_create"`
	TSUpdate time.Time `json:"ts_update"`
	Removed  bool      `json:"removed"`
	Info     H         `json:"info"`
}

// Dump to bytes
func (item BasicFields) Dump() []byte {
	raw, err := sonic.ConfigStd.Marshal(item)
	if err != nil {
		E("json unmarshal error", err)
	}
	return raw
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

func init() {
	sort.Strings(basicFieldsInlineFields)
}
