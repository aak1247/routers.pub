package utils

import (
	"fmt"
	"net/url"
)

func UrlEncode(params map[string]interface{}) string {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}
