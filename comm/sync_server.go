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
			command, err := ReadCommand(conn, conn.RemoteAddr().String())
			if err != nil {
				fmt.Printf("line 25 err: %v", err)
				return err
			}
			err = Reply(conn, command)
			if err != nil {
				fmt.Printf("line 29 err: %v", err)
				return err
			}
		}
	}
}
