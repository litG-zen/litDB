package main

import (
	"fmt"
	"github/litG-zen/litDB/comm"
)

func main() {
	err := comm.RunAsyncServer()
	if err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}
