package core

import (
	"bytes"
	"github/litG-zen/litDB/conf"
	"github/litG-zen/litDB/parser"
	"io"
	"strconv"
	"time"
)

var NIL_RESPONSE []byte = []byte("$-1\r\n")

func evalPING(args []string) []byte {
	var b []byte
	if len(args) > 1 {
		return parser.Encode("Err wrong number of arguments for 'PING' command", false)
	}
	if len(args) == 0 {
		b = parser.Encode("PONG", true)
	} else {
		b = parser.Encode(args[0], false)
	}
	return b
}

func evalSET(args []string) []byte {
	var b []byte
	if len(args) <= 1 {
		return parser.Encode("(error) ERR wrong number of arguments for 'set' command", false)
	}

	var key, value string
	var exDurationMs int64 = -1

	key, value = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			if i+1 >= len(args) {
				return parser.Encode("(error) syntax error", false)
			}
			exDurationSec, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return parser.Encode("(error) ERR value is not an integer or out of range", false)
			}
			i++
			exDurationMs = exDurationSec * 1000
		default:
			return parser.Encode("(error) syntax error", false)
		}
	}

	Put(key, NewObj(value, exDurationMs))
	b = []byte("+OK\r\n")
	return b
}

func evalGET(args []string) []byte {
	var b []byte

	if len(args) != 1 {
		return parser.Encode("(error) ERR wrong number of arguments for 'get' command", false)
	}
	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		b = NIL_RESPONSE
		return b
	}
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		b = NIL_RESPONSE
		return b
	}

	b = parser.Encode(obj.Value, false)
	return b
}

func evalTTL(args []string) []byte {
	var b []byte
	if len(args) != 1 {
		return parser.Encode("error) ERR wrong number of arguments for 'ttl' command", false)
	}

	var key string = args[0]
	obj := store[key]
	if obj == nil {
		b = []byte(":-2\r\n")
		return b
	}

	if obj.ExpiresAt == -1 {
		b = []byte(":-1\r\n")
		return b
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()
	if durationMs < 0 {
		b = []byte(":-2\r\n")
		return b
	}
	b = parser.Encode(int64(durationMs/1000), false)
	return b

}

func evalDEL(args []string) []byte {
	var deletedKeys int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			deletedKeys++
		}
	}

	return parser.Encode(deletedKeys, false)
}

func evalExpire(args []string) []byte {
	if len(args) <= 1 {
		return parser.Encode("(error) ERR wrong number of arguments for 'expire' command", false)
	}

	var key string = args[0]
	expireDuration, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return parser.Encode("(error) ERR value is not an integer or out of range", false)
	}

	obj := Get(key)
	if obj == nil {
		return parser.Encode(":0\r\n", false)
	}

	obj.ExpiresAt = time.Now().UnixMilli() + expireDuration*1000
	return parser.Encode(":1\r\n", false)
}

func EvalAndRespond(c io.ReadWriter, cmds conf.RedisCmds) error {
	var response []byte
	buf := bytes.NewBuffer(response)
	for _, cmd := range cmds {
		switch cmd.Cmd {
		case "PING":
			buf.Write(evalPING(cmd.Args))
		case "GET":
			buf.Write(evalGET(cmd.Args))
		case "SET":
			buf.Write(evalSET(cmd.Args))
		case "TTL":
			buf.Write(evalTTL(cmd.Args))
		case "DEL":
			buf.Write(evalDEL(cmd.Args))
		case "EXPIRE":
			buf.Write(evalExpire(cmd.Args))
		default:
			buf.Write(evalPING(cmd.Args))
		}
	}
	c.Write(buf.Bytes())
	return nil
}
