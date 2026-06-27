package conf

const (
	HOST = "127.0.0.1"
	PORT = 8000
)

type RedisCmd struct {
	Cmd  string
	Args []string
}
