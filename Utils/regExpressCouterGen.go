package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	filename := "CounterTemplate.txt"
	expression, count := GetVeryNewRegExp(filename)
	fmt.Println("COunter Expression ", expression)
	fmt.Println("coiunt ", count)

	f, _ := os.Open(filename)
	data, _ := ioutil.ReadAll(f)
	sdata := strings.TrimSpace((string(data)))

	if mymap, err := getKeyValue(sdata, expression, count); err == nil {
		for key, value := range mymap {
			fmt.Printf("Key:%s,value:%s\n", key, value)
		}
	}

	fmt.Println("######################################")
	fmt.Println("######################################")
	fmt.Println("######################################")
	filename = "StatsTemplate.txt"
	expression = GetNewRegExp(filename)
	fmt.Println("Stats Expression ", expression)

	f, _ = os.Open(filename)
	data, _ = ioutil.ReadAll(f)
	sdata = strings.TrimSpace((string(data)))
	re, _ := regexp.Compile(expression)
	sliceString := re.FindStringSubmatch(sdata)
	fmt.Println("Stats Expression count", len(sliceString))
	if len(sliceString) == 0 {
		return
	}

	mymap := make(map[string]string)
	for i := 1; i < len(sliceString)-1; i += 2 {
		mymap[sliceString[i]] = sliceString[i+1]
	}
	fmt.Printf("Slice is %#v\n", mymap)
	fmt.Println("######################################")
	fmt.Println("######################################")
	fmt.Println("######################################")
	filename = "CfgTemplate.txt"
	expression = GetNewRegExp(filename)
	fmt.Println("CfgTemplate expression", expression)
	f, _ = os.Open(filename)
	data, _ = ioutil.ReadAll(f)
	sdata = strings.TrimSpace((string(data)))

	re, _ = regexp.Compile(expression)
	sliceString = re.FindStringSubmatch(sdata)
	if len(sliceString) == 0 {
		return
	}
	for i := 1; i < len(sliceString)-1; i += 2 {
		mymap[sliceString[i]] = sliceString[i+1]
	}
	fmt.Printf("Slice is %#v\n", mymap)

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

// get Key values as Counter names
func GetVeryNewRegExp(filename string) (string, int) {
	f, _ := os.Open(filename)
	re := regexp.MustCompile(`[\w\s]+\s*:\s*\d+`)
	scanner := bufio.NewScanner(f)
	fmt.Println("filename", filename)
	newstr := ""
	count := 0
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		tt := re.ReplaceAllString(ln, `([\w\s]+):\s*(\d+)`)
		newstr = newstr + `\s*` + tt
		count = count + 1
	}

	return strings.TrimSpace(newstr), count
}

// get Key values as Counter names
func GetNewRegExp(filename string) string {
	f, _ := os.Open(filename)
	re := regexp.MustCompile(`[\w\s]+\s*:\s*\d+`)
	scanner := bufio.NewScanner(f)
	fmt.Println("filename", filename)
	newstr := ""
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		fmt.Println("Scannerl Text", ln)
		tt := re.ReplaceAllString(ln, `([\w\s]+):\s*(\d+)`)
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
