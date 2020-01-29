package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// Place HOlder Functions
func CounterFun(line lineData) {

	filename := line.getFileName()
	text := line.getText()

	// End Delimiter to know that paragraph ended. Not alway new line
	// is deliimter - check the ncdutil output to know what is hte
	// delimiter
	if text == "" {
		Debug("CounterFun:", TempCntrInfo[filename])
		counterHandler(filename, TempCntrInfo[filename])
		TempCntrInfo[filename] = ""
	} else {

		TempCntrInfo[filename] = TempCntrInfo[filename] + text
	}

	return
}

func counterHandler(name string, data string) {

	filename := "CounterTemplate.txt"
	Debug("Raju: CounterFun:", name, data)
	// expression := `\s+Port: (\d+)\s+Good UC Pkt Rcvd:\s+(\d+)`
	expression, keyslice := GetRegExp(filename)
	r, _ := regexp.Compile(expression)
	if result := r.FindAllStringSubmatch(data, -1); len(result) != 0 {
		for _, element := range result {
			Debug("Pkt count port :", element[1])
			Debug("Pkt count :", element[2])
			processCounterValues(name, element, keyslice)
		}
	}
}

func processCounterValues(name string, strslice []string, keyslice []string) {
	Debug("Raju: ProcessCounterValues:", strslice[1])

	// Any of these values should not be zero. Get the line numbers
	// from TextCounter.txt file [set num]
	// var BadArray = []int{5, 6, 16, 21, 22, 23, 24, 26, 27, 28, 29, 30}
	var BadArray = []int{3, 4, 14, 19, 20, 21, 22, 24, 25, 26, 27, 28}

	for i := 0; i < len(BadArray); i++ {
		index := BadArray[i]
		if intValue, err := strconv.Atoi(strslice[index]); err == nil {
			if intValue != 0 {
				errstr := fmt.Sprintf("ERROR: %s:Port:%s have Error:%s:%d\n",
					name, strslice[1], keyslice[index-1], intValue)
				SendError(errstr)
			}
		} else {
			Error("Count not convert the index Integr", index)
		}
	}

}
