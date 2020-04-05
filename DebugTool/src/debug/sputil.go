package debug

import (
	"fmt"
	"io/ioutil"
	"os"
)

type SputilInfo struct {
	Server string
	// LDFS CORE  IDString key, value
	//			  IDString
	// Mariner    IDstring
	// 			  IDString
	Mymapmap map[string]map[string]map[string]string
}

func ProcessSputilInfo(path string, servername string) {

	f, _ := os.Open(path)

	b, _ := ioutil.ReadAll(f)
	Debug("Sptuils", string(b))

	sputil := SputilInfo{Server: servername}
	fmt.Println("sputils", sputil)

}

func GetSputilInfo() interface{} {
	sputil := SputilInfo{}
	return readFromDb(sputil)
}
