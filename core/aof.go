package core

import (
	"fmt"
	"github/litG-zen/litDB/conf"
	"github/litG-zen/litDB/parser"
	"log"
	"os"
	"strings"
)

func dumpkey(key string, obj *Obj, file *os.File) {
	if obj == nil {
		return
	}
	cmd := fmt.Sprintf("SET %s %s", key, obj.Value)
	tokens := strings.Split(cmd, " ")
	file.Write(parser.Encode(tokens, false))
}

func DumpAllAOF() {
	file, err := os.OpenFile(conf.AOF_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening AOF file: %v\n", err)
		return
	}
	for k, obj := range store {
		dumpkey(k, obj, file)
	}
	log.Printf("AOF rewrite completed. %d keys dumped to %s\n", len(store), conf.AOF_FILE)
}
