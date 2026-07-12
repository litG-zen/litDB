package parser

import (
	"bytes"
	"fmt"
)

// RESP specification: https://redis.io/docs/reference/protocol-spec/

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte("+" + v + "\r\n")
		}
		return encodeString(v)

	case int64, int, int8, int32, int16:
		return []byte(fmt.Sprintf(":%d\r\n", v))

	case []string:
		var b []byte
		buff := bytes.NewBuffer(b)
		for _, b := range value.([]string) {
			buff.Write(encodeString(b))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(v), buff.Bytes()))

	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v.Error()))

	default:
		return []byte{}
	}
}
