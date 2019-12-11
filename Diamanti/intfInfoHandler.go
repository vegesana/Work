package main

import (
	"regexp"
)

func IntfInfoFun(line lineData) {

	filename := line.getFileName()
	text := line.getText()

	// End Delimiter to know that paragraph ended. Not alway new line
	// is deliimter - check the ncdutil output to know what is hte
	// delimiter
	if text == "" {
		Debug("Raju:", TempIntfInfo[filename])
		intfInfoHandler(filename, TempIntfInfo[filename])
		TempIntfInfo[filename] = ""
	} else {
		TempIntfInfo[filename] = TempIntfInfo[filename] + text
	}

	return
}

func intfInfoHandler(name string, data string) {

	r, _ := regexp.Compile(`Port: (\d+), Link: \w+, Pvid: \d+, MTU: \d+, Lane0_Tx_Polarity:\s+` +
		`\w+, Lane0_Rx_Polarity: \w+\s+Protected_Mode: \w+, dstPort: (\d+), dstHwDev: \d+`)
	if result := r.FindAllStringSubmatch(data, -1); len(result) != 0 {
		for _, element := range result {
			// we need this 2 times to trip two extensions
			sport := GPort{name, element[1]}
			dport := GPort{name, element[2]}
			pininfo := PinInfo{sport, dport}
			Error("RajuPost", pininfo)
			Gdata.GWriteCh <- pininfo
		}
	}
}
