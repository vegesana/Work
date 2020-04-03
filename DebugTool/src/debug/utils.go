package debug

import (
	"bufio"
	"errors"
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
	fmt.Println("filename", filename)
	newstr := ""
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		fmt.Println("Scannerl Text", ln)
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

func getKeyValue(data, express string, count int) (map[string]string, error) {

	var err = errors.New("Error")

	r, _ := regexp.Compile(express)
	mymap := make(map[string]string)
	if result := r.FindAllStringSubmatch(data, -1); len(result) != 0 {
		err = nil
		for _, element := range result {
			for i := 0; i < count; i++ {
				j := 2*i + 1
				mymap[element[j]] = element[j+1]
			}
		}
	}
	return mymap, err

}

func getNewKeyValue(data, express string) (map[int]map[string]string, error) {

	var err = errors.New("Error")

	r, _ := regexp.Compile(express)
	mymap := make(map[int]map[string]string)

	// Result is slice of slices
	if resultSliceSlice := r.FindAllStringSubmatch(data, -1); len(resultSliceSlice) != 0 {
		for index, slice := range resultSliceSlice {
			mymap[index] = make(map[string]string)
			mylen := len(slice)
			for j := 0; j < (mylen-1)/2; j++ {
				k := 2*j + 1
				//fmt.Printf("Index:%d,key:%svalue:%s\n", index, slice[k], slice[k+1])
				mymap[index][slice[k]] = slice[k+1]
				err = nil
			}
		}
	}
	return mymap, err

}
