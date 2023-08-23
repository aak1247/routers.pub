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

// 扁平化多级map，key使用keySep连接
func FlattenMap(m map[string]interface{}, keySep string, keyPrefix string) map[string]interface{} {
	flattenedMap := make(map[string]interface{})
	for k, v := range m {
		if keyPrefix != "" {
			k = keyPrefix + keySep + k
		}
		flattenedMap[k] = v // 先把当前层级的key-value放入结果
		if vMap, ok := v.(map[string]interface{}); ok {
			for fk, fv := range FlattenMap(vMap, keySep, k) {
				flattenedMap[fk] = fv
			}
		}
	}
	return flattenedMap
}

func GetFuzzyMap(m map[string]interface{}, key string) map[string]interface{} {
	fuzzyMap := make(map[string]interface{})
	for k, v := range m {
		if strings.Contains(k, key) {
			fuzzyMap[k] = v
		}
	}
	return fuzzyMap
}

func TraverseMapString(m map[string]interface{}, f func(key string, value string) string) {
	for k, v := range m {
		switch v.(type) {
		case string:
			m[k] = f(k, v.(string))
		case map[string]interface{}:
			TraverseMapString(v.(map[string]interface{}), f)
		case []interface{}:
			for _, v := range v.([]interface{}) {
				switch v.(type) {
				case string:
					m[k] = f(k, v.(string))
				case map[string]interface{}:
					TraverseMapString(v.(map[string]interface{}), f)
				}
			}
		case map[string]string:
			for k, v := range v.(map[string]string) {
				m[k] = f(k, v)
			}
		}
	}
}

func TraverseStringMap(m map[string]string, f func(key string, value string)) {
	for k, v := range m {
		f(k, v)
	}
}
