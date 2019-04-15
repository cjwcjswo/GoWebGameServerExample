package util

import "encoding/json"

type Parser interface {
	Decoding(data []byte, element interface{}) error
	Encoding(data interface{}) ([]byte, error)
}

type JsonParser struct{}

func (parser JsonParser) Decoding(data []byte, element interface{}) error {
	return json.Unmarshal(data, element)
}

func (parser JsonParser) Encoding(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}
