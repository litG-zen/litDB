package comm

import (
	"fmt"
	"github/litG-zen/litDB/conf"
	"net"
)

func RunSyncServer() error {
	listner, err := net.Listen("tcp", conf.HOST+":"+fmt.Sprint(conf.PORT))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Listening on port " + fmt.Sprint(conf.PORT))

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Printf("line 20 err: %v", err)
			return err
		}
		for {
			command, err := ReadCommands(conn, conn.RemoteAddr().String())
			if err != nil {
				return err
			}
			err = Reply(conn, command)
			if err != nil {
				return err
			}
		}
	}
}
