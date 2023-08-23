package utils

import "strings"

func CamelToSnake(input string) string {
	var output = ""
	for i := 0; i < len(input); i += 1 {
		if i > 0 && input[i] >= 'A' && input[i] <= 'Z' {
			output += "_"
		}
		output += string(input[i])
	}
	return strings.ToLower(output)
}
