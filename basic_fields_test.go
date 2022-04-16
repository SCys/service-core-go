package core

import (
	"testing"
)

type TestBasicFieldsA struct {
	BasicFields

	A string `json:"a"`
	B int    `json:"b"`
}

func (t TestBasicFieldsA) TableName() string {
	return "t_abc"
}

func TestBasicFieldsGenFieldsAndArgumentsA(t *testing.T) {
	item := TestBasicFieldsA{
		BasicFields: NewBasicFields(),

		A: "abc",
		B: 123,
	}

	fields, arguments := basicFieldsGenFieldsAndArguments(item)
	if len(fields) != 7 {
		t.Errorf("Expected 7 fields, got %d", len(fields))
	}

	if len(arguments) != 7 {
		t.Errorf("Expected 7 arguments, got %d", len(arguments))
	}

	if fields[5] != "abc" {
		t.Errorf("Expected field 5 to be 'abc', got '%s'", fields[5])
	}
	if fields[6] != 123 {
		t.Errorf("Expected field 6 to be 123, got '%d'", fields[6])
	}

	t.Failed()
}

func TestBasicFieldsGenFieldsAndArgumentsB(t *testing.T) {
	item := TestBasicFieldsA{
		BasicFields: NewBasicFields(),

		A: "abc",
		B: 123,
	}

	fields, arguments := basicFieldsGenFieldsAndArguments(&item)
	if len(fields) != 7 {
		t.Errorf("Expected 7 fields, got %d", len(fields))
	}

	if len(arguments) != 7 {
		t.Errorf("Expected 7 arguments, got %d", len(arguments))
	}

	if fields[5] != "abc" {
		t.Errorf("Expected field 5 to be 'abc', got '%s'", fields[5])
	}
	if fields[6] != 123 {
		t.Errorf("Expected field 6 to be 123, got '%d'", fields[6])
	}

	t.Failed()
}
