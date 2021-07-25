package core

import (
	"github.com/valyala/fastjson"
)

var parserPool = fastjson.ParserPool{}

// ParseJSON parse raw to fastjson.Value
func ParseJSON(raw string) (*fastjson.Value, error) {
	p := ParserGet()

	v, e := p.Parse(raw)

	ParserPut(p)

	return v, e
}

// ParseJSONBytes parse json bytes
func ParseJSONBytes(raw []byte) (*fastjson.Value, error) {
	return ParseJSON(String(raw))
}

// ParserGet from pool
func ParserGet() *fastjson.Parser {
	return parserPool.Get()
}

// ParserPut return to pool
func ParserPut(p *fastjson.Parser) {
	parserPool.Put(p)
}
