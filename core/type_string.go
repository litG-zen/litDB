package core

import "strconv"

func deduceTypeEncoding(v string) (uint8, uint8) {
	oType := OBJ_TYPE_STRING

	if _, err := strconv.ParseInt(v, 10, 64); err == nil { // if we  are able to convert the object into Int, then the obj type  is string.
		return oType, OBJ_ENCODING_INT
	}
	if len(v) <= 44 {
		return oType, OBJ_ENCODING_EMBSTR
	}
	return oType, OBJ_ENCODING_RAW
}
