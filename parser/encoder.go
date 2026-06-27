package parser

import "fmt"

// RESP specification: https://redis.io/docs/reference/protocol-spec/

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte("+" + v + "\r\n")
		}
		return []byte("$" + fmt.Sprintf("%d", len(v)) + "\r\n" + v + "\r\n")

	case int64, int, int8, int32, int16:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	}

	return []byte{}
}
