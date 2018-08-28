package utils

import (
	"encoding/json"
	"fmt"
)

// FormatResult 格式化http返回参数
func FormatResult(v interface{}) []byte {
	result := &struct {
		Code byte        `json:"code"`
		Data interface{} `json:"data"`
	}{
		Code: 0,
		Data: v,
	}
	message, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json.Marshal error.")
	}
	return message
}
