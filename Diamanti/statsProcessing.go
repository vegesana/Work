package main

import (
	"fmt"
	"regexp"
	"strings"
)

func StatsInfoFun(line lineData) {

	filename := line.getFileName()
	text := strings.TrimSpace(line.getText())

	// End Delimiter to know that paragraph ended. Not alway new line
	// is deliimter - check the ncdutil output to know what is hte
	// delimiter
	if text == "" {
		Debug("StatsInfoFunc:", TempStatsInfo[filename])
		statsHandler(filename, TempStatsInfo[filename])
		TempStatsInfo[filename] = ""
	} else {

		TempStatsInfo[filename] = TempStatsInfo[filename] + text
	}

	return
}

func statsHandler(name string, data string) {
	mymap := map[string]string{}
	filename := "StatsTemplate.txt"

	Debug("Raju: StatsTemplate:", name, data)
	expression := GetNewRegExp(filename)
	re, _ := regexp.Compile(expression)

	sliceString := re.FindStringSubmatch(data)

	if len(sliceString) == 0 {
		Debug("Raju: stathandler NO expression match ", sliceString)
		return
	}

	for i := 1; i < len(sliceString)-1; i += 2 {
		mymap[sliceString[i]] = sliceString[i+1]
	}

	processStatsValues(name, mymap)
}
func processStatsValues(name string, mymap map[string]string) {
	var errstr string
	var mystr string
	Debug("Raju: ProcessStatsValues:", mymap)

	mymapint := convertMapStringToMapInt(mymap)

	mystr = "nm_hello"
	if mymapint[mystr] != 1 {
		errstr = fmt.Sprintf("%s:%s is Invalid value: %d\n",
			name, mystr, mymapint[mystr])
		SendError(errstr)
	}

	mystr = "nm_disconnect"
	if mymapint[mystr] != 0 {
		errstr = fmt.Sprintf("%s:%s is Invalid value: %d\n",
			name, mystr, mymapint[mystr])
		SendError(errstr)
	}
}
