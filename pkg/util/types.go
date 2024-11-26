package util

import (
	"reflect"
	"strings"
)

// Removes package from result
// Removes pointer from result
func GetBaseType(myvar any) string {
	result := GetType(myvar)
	if i := strings.Index(result, "."); i != -1 {
		result = result[i+1:]
	}
	if strings.HasPrefix("*", result) {
		result = result[1:]
	}
	return result
}

func GetType(myvar any) string {
	result := reflect.TypeOf(myvar).String()
	return result
}
