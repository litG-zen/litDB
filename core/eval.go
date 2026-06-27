package core

import (
	"errors"
	"github/litG-zen/litDB/conf"
	"github/litG-zen/litDB/parser"
	"io"
	"strconv"
	"time"
)

var NIL_RESPONSE []byte = []byte("$-1\r\n")

func evalPING(c io.ReadWriter, args []string) error {
	var b []byte
	if len(args) > 1 {
		return errors.New("Err wrong number of arguments for 'PING' command")
	}
	if len(args) == 0 {
		b = parser.Encode("PONG", true)
	} else {
		b = parser.Encode(args[0], false)
	}
	_, err := c.Write(b)
	return err
}

func evalSET(c io.ReadWriter, args []string) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'set' command")
	}

	var key, value string
	var exDurationMs int64 = -1

	key, value = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			if i+1 >= len(args) {
				return errors.New("(error) syntax error")
			}
			exDurationSec, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return errors.New("(error) ERR value is not an integer or out of range")
			}
			i++
			exDurationMs = exDurationSec * 1000
		default:
			return errors.New("(error) syntax error")
		}
	}

	Put(key, NewObj(value, exDurationMs))
	c.Write([]byte("+OK\r\n"))
	return nil
}

func evalGET(c io.ReadWriter, args []string) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'get' command")
	}
	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		c.Write(NIL_RESPONSE)
		return nil
	}
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		c.Write(NIL_RESPONSE)
		return nil
	}

	c.Write(parser.Encode(obj.Value, false))
	return nil
}

func evalTTL(c io.ReadWriter, args []string) error {
	if len(args) != 1 {
		return errors.New("error) ERR wrong number of arguments for 'ttl' command")
	}

	var key string = args[0]
	obj := store[key]
	if obj == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	if obj.ExpiresAt == -1 {
		c.Write([]byte(":-1\r\n"))
		return nil
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()
	if durationMs < 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	c.Write(parser.Encode(int64(durationMs/1000), false))
	return nil

}

func evalDEL(c io.ReadWriter, args []string) error {
	var deletedKeys int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			deletedKeys++
		}
	}

	c.Write(parser.Encode(deletedKeys, false))
	return nil
}

func evalExpire(c io.ReadWriter, args []string) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'expire' command")
	}

	var key string = args[0]
	expireDuration, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("(error) ERR value is not an integer or out of range")
	}

	obj := Get(key)
	if obj == nil {
		c.Write([]byte(":0\r\n"))
	}

	obj.ExpiresAt = time.Now().UnixMilli() + expireDuration*1000
	c.Write([]byte(":1\r\n"))
	return nil
}

func EvalAndRespond(c io.ReadWriter, cmd *conf.RedisCmd) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(c, cmd.Args)
	case "GET":
		return evalGET(c, cmd.Args)
	case "SET":
		return evalSET(c, cmd.Args)
	case "TTL":
		return evalTTL(c, cmd.Args)
	case "DEL":
		return evalDEL(c, cmd.Args)
	case "EXPIRE":
		return evalExpire(c, cmd.Args)
	default:
		return evalPING(c, cmd.Args)
	}
}
