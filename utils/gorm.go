package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type StringSlice []string

func (ss *StringSlice) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// Scan 实现方法
func (ss *StringSlice) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &ss)
}

type StringMap map[string]string

func (sm *StringMap) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

// Scan 实现方法
func (sm *StringMap) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &sm)
}

func (sm *StringMap) Set(key, value string) {
	if *sm == nil {
		*sm = make(map[string]string)
	}
	(*sm)[key] = value
}
