package main

import (
	"fmt"
	"regexp"
	"strings"
)

func CfgInfoFun(line lineData) {

	filename := line.getFileName()
	text := strings.TrimSpace(line.getText())

	cfgHandler(filename, text)

	return
}

func cfgHandler(name string, data string) {
	mymap := map[string]string{}
	filename := "CfgTemplate.txt"

	Debug("Raju: CfgTemplate:", name, data)
	expression := GetNewRegExp(filename)
	Debug("Raju: reexp:", expression)
	re, _ := regexp.Compile(expression)

	sliceString := re.FindStringSubmatch(data)

	if len(sliceString) == 0 {
		Debug("Raju: no expression match ", sliceString)
		return
	}
	Debug("Raju: cfgHandler slicestring:", sliceString)
	for i := 1; i < len(sliceString)-1; i += 2 {
		mymap[sliceString[i]] = sliceString[i+1]
	}

	processCfgValues(name, mymap)
}
func processCfgValues(name string, mymap map[string]string) {
	var errstr string
	var mystr string
	Debug("Raju: CfgStatsValues:", mymap)

	mymapint := convertMapStringToMapInt(mymap)

	mystr = "hb_state"
	if mymapint[mystr] != 1 {
		errstr = fmt.Sprintf("%s:%s is No heart beat : %d\n",
			name, mystr, mymapint[mystr])
		SendError(errstr)
	}

	mystr = "curr_uptime"
	cmystr := "last_uptime"
	if mymapint[mystr]-mymapint[cmystr] > 90 {
		errstr = fmt.Sprintf("%s:%s Heart beat > 90: %d\n",
			name, mystr, mymapint[mystr])
		SendError(errstr)
	}

}
