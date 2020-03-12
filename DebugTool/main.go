package main

import (
	"DebugTool/src/rest"
	"fmt"
)

func main() {
	fmt.Println("vim-go")
	robj := rest.Init()
	robj.Start()

	for {
	}
}
