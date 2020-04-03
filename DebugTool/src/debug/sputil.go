package debug

import (
	"io/ioutil"
	"os"
)

type SputilInfo struct {
	name     string
	mymapmap map[int]map[string]string
}

func ProcessSputilInfo(path string, servername string) {

	f, _ := os.Open(path)

	b, _ := ioutil.ReadAll(f)
	Debug("Sptuils", string(b))

}

func GetSputilInfo() interface{} {
	sputil := SputilInfo{}
	return readFromDb(sputil)
}
