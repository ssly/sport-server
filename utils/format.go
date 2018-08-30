package utils

import (
	"encoding/json"
	"fmt"
)

// FormatResult 格式化http返回参数
func FormatResult(code int, v interface{}) []byte {
	result := &struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}{
		Code: code,
		Data: v,
	}
	message, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json.Marshal error.")
	}
	return message
}
