package debug

import (
	"DebugTool/src/utils"
	"fmt"
	"io/ioutil"
	"os"
)

func processBoardInfo(path string, servername string) error {

	if _, ok := SystemMap[servername]; ok {
		fmt.Println("New Path that is aready exist", path)
		return nil
	}

	f, _ := os.Open(path)
	allText, _ := ioutil.ReadAll(f)
	productid := utils.GetValueOfStr(string(allText), "Product id", ":")
	board := utils.GetValueOfStr(string(allText), "Board", ":")

	// this should send write on to the channel to update data
	writeToDb(SystemInfo{board, productid, servername})

	return nil
}
