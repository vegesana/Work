package debug

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Info(param ...interface{}) {
	fmt.Println(param)
}
func Error(param ...interface{}) {
	fmt.Println(param)
}
func Input(param ...interface{}) {
	fmt.Println(param)
}

func Dump(param ...interface{}) {
	fmt.Println(param)
}

func Debug(param ...interface{}) {
	//fmt.Println(param)
	return
}
func DumpErrors() {
	fmt.Println("ERRORS")
	fmt.Println("==================")
	for _, err := range MyErrSlice {
		fmt.Println(err.MyErr)
	}
	fmt.Println("==================")
}

func GetRegExp(filename string) (string, []string) {
	f, _ := os.Open(filename)

	myslice := []string{}
	newstr := ""
	re := regexp.MustCompile(`:\s*\d+`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		newslice := strings.Split(ln, ":")
		myslice = append(myslice, newslice[0])
		tt := re.ReplaceAllString(ln, `:\s*(\d+)`)
		newstr = newstr + `\s+` + tt
	}
	return newstr, myslice
}

// get Key values as Counter names
func GetNewRegExp(filename string) string {
	f, _ := os.Open(filename)
	re := regexp.MustCompile(`\w+\s*:\s*\d+`)
	scanner := bufio.NewScanner(f)
	newstr := ""
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		tt := re.ReplaceAllString(ln, `(\w+):\s*(\d+)`)
		newstr = newstr + `\s*` + tt
	}

	return newstr
}

func getValueOfStr(ln string, keystr string, delimit string) string {
	re, _ := regexp.Compile(`\s*` + keystr + `\s*` + delimit + `\s*(\w+.*)`)
	value := re.FindStringSubmatch(ln)
	if len(value) >= 2 {
		return value[1]
	}
	return ""
}

func convertMapStringToMapInt(mymap map[string]string) map[string]int {
	newmap := map[string]int{}

	for key, value := range mymap {
		if intvalue, err := strconv.Atoi(value); err == nil {
			newmap[key] = intvalue
		} else {
			fmt.Println("String converstion error for MAP to MAP ", value)
		}
	}

	return newmap
}
