package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func DumpErrors() {
	fmt.Println("ERRORS")
	fmt.Println("==================")
	for _, err := range MyErrSlice {
		fmt.Println(err.myerr)
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
