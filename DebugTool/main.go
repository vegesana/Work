package main

import (
	"DebugTool/src/rest"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("vim-go")
	robj := rest.Init()
	robj.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	fmt.Println("Exit main")
}
