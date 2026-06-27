package parser

// RESP specification: https://redis.io/docs/reference/protocol-spec/

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}
	ts := value.([]interface{}) //type assertion to []interface{}
	tokens := make([]string, len(ts))
	for i := range ts {
		tokens[i] = ts[i].(string)
	}
	return tokens, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	value, _, err := decodeOne(data)
	return value, err
}

func decodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, nil
	}

	switch data[0] {
	case '+':
		return decodeSimpleString(data)
	case '-':
		return decodeError(data)
	case ':':
		return decodeInteger(data)
	case '$':
		return decodeBulkString(data)
	case '*':
		return decodeArray(data)
	default:
		return nil, 0, nil
	}
}

func decodeSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

func decodeError(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

func decodeInteger(data []byte) (int64, int, error) {
	pos := 1
	var value int64
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil
}

func decodeBulkString(data []byte) (string, int, error) {
	pos := 1
	len, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos : pos+len]), pos + len + 2, nil
}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return length, pos
}

func decodeArray(data []byte) (interface{}, int, error) {
	pos := 1
	count, delta := readLength(data[pos:])
	pos += delta

	var elems []interface{} = make([]interface{}, count)
	for i := range elems {
		elem, delta, err := decodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}
