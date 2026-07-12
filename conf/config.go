package conf

const (
	HOST      = "127.0.0.1"
	PORT      = 8000
	KEYSLIMIT = 10
)

type RedisCmd struct {
	Cmd  string
	Args []string
}

type RedisCmds []*RedisCmd

const AOF_FILE = "./litdb.aof"
