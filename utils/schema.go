package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type JSType string

const (
	JSTypeString  JSType = "string"
	JSTypeNumber  JSType = "number"
	JSTypeBoolean JSType = "boolean"
	JSTypeObject  JSType = "object"
	JSTypeArray   JSType = "array"
)

type JSONSchema struct {
	// 类型
	Type       JSType                `json:"type"`
	Default    interface{}           `json:"default"`
	Properties map[string]JSONSchema `json:"properties"`
}

type JSONSchemaValue struct {
	// 类型
	Type       JSType                `json:"type"`
	Value      interface{}           `json:"value"`
	Properties map[string]JSONSchema `json:"properties"`
}

func NewJSONSchema(s string) *JSONSchema {
	js := new(JSONSchema)
	json.Unmarshal([]byte(s), js)
	return js
}

func (js *JSONSchema) AddProperty(name string, schema JSONSchema) {
	if js.Properties == nil {
		js.Properties = make(map[string]JSONSchema)
	}
	js.Properties[name] = schema
}

func (js *JSONSchema) AddPropertyByType(name string, t JSType) {
	js.AddProperty(name, JSONSchema{
		Type: t,
	})
}

func (js *JSONSchema) ParseData(data map[string]interface{}) *JSONSchemaValue {
	jsv := new(JSONSchemaValue)
	jsv.Type = js.Type
	jsv.Properties = make(map[string]JSONSchema)
	return jsv
}

func (js *JSONSchema) GetDefaultValue() interface{} {
	if js == nil {
		return nil
	}
	if js.Default == nil {
		switch js.Type {
		case JSTypeString:
			return ""
		case JSTypeNumber:
			return 0
		case JSTypeBoolean:
			return false
		case JSTypeObject:
			m := make(map[string]interface{})
			for k, v := range js.Properties {
				m[k] = v.GetDefaultValue()
			}
			return m
		case JSTypeArray:
			a := make([]interface{}, 0)
			for _, v := range js.Properties {
				a = append(a, v.GetDefaultValue())
			}
			return a
		}
	}
	return js.Default
}

func (js *JSONSchema) Value() (driver.Value, error) {
	return json.Marshal(js)
}

// Scan 实现方法
func (js *JSONSchema) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &js)
}
