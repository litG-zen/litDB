package comm

import (
	"fmt"
	"github/litG-zen/litDB/conf"
	"github/litG-zen/litDB/core"
	"github/litG-zen/litDB/parser"
	"io"
	"strings"
	"syscall"
)

// //////////////////////////////////////////////////////////////////////////////////

type FDComm struct {
	Fd         int
	ClientAddr string
}

func (c *FDComm) Read(b []byte) (int, error) {
	return syscall.Read(c.Fd, b)
}

func (c *FDComm) Write(b []byte) (int, error) {
	return syscall.Write(c.Fd, b)
}

// //////////////////////////////////////////////////////////////////////////////////

func toArrayString(tokens []interface{}) ([]string, error) {
	as := make([]string, len(tokens))
	for i := range tokens {
		as[i] = tokens[i].(string)
	}
	return as, nil
}

// Socket read and write functions
func ReadCommands(c io.ReadWriter, clientAddr string) (conf.RedisCmds, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	values, err := parser.Decode(buf[:n])
	if err != nil {
		return nil, err
	}
	var cmds []*conf.RedisCmd = make([]*conf.RedisCmd, 0)
	for _, value := range values {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmd := &conf.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func respondError(err error, c io.Writer) {
	errorMessage := fmt.Sprintf("-%s\r\n", err)
	c.Write([]byte(errorMessage))
}

func Reply(c io.ReadWriter, response conf.RedisCmds) error {
	err := core.EvalAndRespond(c, response)
	if err != nil {
		respondError(err, c)
	}
	return nil
}
