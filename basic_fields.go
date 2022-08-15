package core

import (
	"sort"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rs/xid"
	"github.com/valyala/fastjson"
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
	raw, err := sonic.Marshal(item)
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

// LoadBasicFields load with default fields
func LoadBasicFields(i any, value *fastjson.Value) {
	x, ok := i.(*BasicFields)
	if !ok {
		return
	}

	x.ID = String(value.GetStringBytes("id"))
	x.TSCreate, _ = time.Parse(time.RFC3339Nano, String(value.GetStringBytes("ts_create")))
	x.TSUpdate, _ = time.Parse(time.RFC3339Nano, String(value.GetStringBytes("ts_update")))
	x.Removed = value.GetBool("removed")

	_ = sonic.Unmarshal(value.GetStringBytes("info"), &x.Info)
}

func init() {
	sort.Strings(basicFieldsInlineFields)
}
