package utils

import "strings"

// 对json结构的多级map 获取多级路径指定的值
func DepGetFromMap(m map[string]interface{}, key string, keySep string) interface{} {
	if m == nil {
		return nil
	}
	// 按照 keySep 分割 key 得到多级路径
	keys := strings.Split(key, keySep)
	// 从 map 中获取值
	var value interface{} = m
	for _, k := range keys {
		if value == nil {
			return nil
		}
		value = value.(map[string]interface{})[k]
	}
	return value
}

// 对json结构的多级map 设置多级路径指定的值
func DepSetToMap(m map[string]interface{}, key string, keySep string, value interface{}) {
	if m == nil {
		return
	}
	// 按照 keySep 分割 key 得到多级路径
	keys := strings.Split(key, keySep)
	// 从 map 中获取值
	var v interface{} = m
	for i, k := range keys {
		if v == nil {
			return
		}
		if i == len(keys)-1 {
			v.(map[string]interface{})[k] = value
		} else {
			v = v.(map[string]interface{})[k]
		}
	}
}
