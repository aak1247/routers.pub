package framework

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

type (
	SearchInfo struct {
		Enabled bool
		Key     string
	}
)

func (si *SearchInfo) ParseTag(str string) *SearchInfo {
	if str != "" && str != "false" && str != "-" {
		si.Enabled = true
	}
	tags := strings.Split(str, ";")
	for _, tag := range tags {
		if strings.Contains(tag, "key:") {
			si.Key = strings.Split(tag, ":")[1]
		}
	}
	return si
}

func BindQuery(c *gin.Context, s interface{}) error {
	// 反射
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < t.NumField(); i++ {
		si := (&SearchInfo{}).ParseTag(t.Field(i).Tag.Get("search"))
		if si.Enabled {
			if si.Key != "" {
				v.Field(i).SetString(c.Query(si.Key))
			} else {
				v.Field(i).SetString(c.Query(t.Field(i).Name))
			}
		}
	}
	return nil
}

func GetSearchKey(s interface{}, prefix string) string {
	// 遍历所有属性，拼接成字符串
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	// 反射
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < t.NumField(); i++ {
		si := (&SearchInfo{}).ParseTag(t.Field(i).Tag.Get("search"))
		if !v.Field(i).IsZero() && si.Enabled {
			buffer.WriteString(fmt.Sprintf("%v.", v.Field(i)))
		}
	}
	return buffer.String()
}
