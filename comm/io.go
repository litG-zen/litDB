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

// Socket read and write functions
func ReadCommand(c io.ReadWriter, clientAddr string) (*conf.RedisCmd, error) {
	buf := make([]byte, 1024)

	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, io.EOF
	}

	tokens, err := parser.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to decode command: %v", err)
	}
	return &conf.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c io.Writer) {
	errorMessage := fmt.Sprintf("-%s\r\n", err)
	c.Write([]byte(errorMessage))
}

func Reply(c io.ReadWriter, response *conf.RedisCmd) error {
	err := core.EvalAndRespond(c, response)
	if err != nil {
		respondError(err, c)
	}
	return nil
}
