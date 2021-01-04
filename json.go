package core

import (
	"github.com/valyala/fastjson"
)

var parserPool fastjson.ParserPool

// ParseJSON parse json string
func ParseJSON(raw string) (*fastjson.Value, error) {
	p := ParserGet()

	v, e := p.Parse(raw)

	ParserPut(p)

	return v, e
}

// ParseJSONBytes parse json bytes
func ParseJSONBytes(raw []byte) (*fastjson.Value, error) {
	p := ParserGet()

	v, e := p.Parse(String(raw))

	ParserPut(p)

	return v, e
}

// ParserGet from pool
func ParserGet() *fastjson.Parser {
	return parserPool.Get()
}

// ParserPut return to pool
func ParserPut(p *fastjson.Parser) {
	parserPool.Put(p)
}

func init() {
	parserPool = fastjson.ParserPool{}
}
