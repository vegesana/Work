package debug

import (
	"fmt"
	"io/ioutil"
	"os"
)

type CtrlInfo struct {
	name     string
	mymapmap map[int]map[string]string
}

func ProcessCtrlInfo(path string, servername string) {

	f, _ := os.Open(path)
	b, _ := ioutil.ReadAll(f)

	expression := autGeneratedNVMEoEUtilExpr()
	if mymap, err := getNewKeyValue(string(b), expression); err == nil {
		processCtrlValuesNew(servername, mymap)
	} else {
		fmt.Println("EROR Processing nvmeoeutil ***** ")
	}
}
func processCtrlValuesNew(name string, mymap map[int]map[string]string) {
	newmap := map[int]map[string]string{}

	for index, map1 := range mymap {
		if map1["allocated"] == "1" {
			if _, ok := newmap[index]; !ok {
				newmap[index] = make(map[string]string)
			}
			newmap[index] = map1
		}
	}
	writeToDb(CtrlInfo{name, newmap})
}
func GetCtrlInfo() interface{} {
	ctrl := CtrlInfo{}
	return readFromDb(ctrl)
}
