package core

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"

	"github.com/rs/xid"
)

// H alias map[string]interface{}
type H = map[string]interface{}

// HNil empty struct
var HNil = H{}

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

// NewBasicFields basicFields init
func NewBasicFields() BasicFields {
	now := Now()

	return BasicFields{
		ID:       xid.New().String(),
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
