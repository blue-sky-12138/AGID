package qjson

// 本库是快速编码解码json的第三方库json-iterator/go的封装

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func MarshalToString(data interface{}) (string, error) {
	bytes, err := Marshal(data)
	return string(bytes), err
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func UnmarshalString(data string, v interface{}) error {
	return Unmarshal([]byte(data), v)
}
