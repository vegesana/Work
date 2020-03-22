package debug

import (
	"DebugTool/src/utils"
	"fmt"
	"io/ioutil"
	"os"
)

func processBoardInfo(path string, servername string) error {

	f, _ := os.Open(path)
	fmt.Println("inside processBoardInfo", servername)
	allText, _ := ioutil.ReadAll(f)
	productid := utils.GetValueOfStr(string(allText), "Product id", ":")
	board := utils.GetValueOfStr(string(allText), "Board", ":")

	// this should send write on to the channel to update data
	fmt.Println("write to processBoardInfo", servername)
	writeToDb(SystemInfo{board, productid, servername})

	fmt.Println("end processBoardInfo", servername)

	return nil
}
