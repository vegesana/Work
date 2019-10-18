package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Gch chan string

var myfile *os.File

type Data struct {
	mych     chan string
	filename string
	funptr   func(string, Data)
}

func Debug(param ...interface{}) {
	fmt.Fprintln(myfile, param)
	fmt.Println(param)
}
func main() {

	// String Slice - This will be adjusted based on the text in the
	// file
	defer myfile.Close()
	BUFSIZE := 40
	myfile, _ = os.OpenFile("debug.txt", os.O_WRONLY|os.O_CREATE, 0666)
	junkData := Data{mych: make(chan string)}
	Gch = junkData.mych
	go goRoutine("Junk Channel", junkData)
	strslice := []string{"Network Interface Info", "NIC Interface Info",
		"NVMEOE Interface Info", "NIF Interface Info", "LIF Interface Info",
		"Statistics Info", "Cfg Info", "CPSS Info", "PCL Info",
		"VLAN Info", "MAC Info"}
	myMap := make(map[string]Data)
	myMap[strslice[10]] = Data{
		funptr:   processMacInfo,
		filename: "MacData.txt",
	}

	for _, str := range strslice {
		if val, ok := myMap[str]; ok {
			val.mych = make(chan string, BUFSIZE)
			myMap[str] = val
		} else {
			myMap[str] = Data{
				mych: make(chan string, BUFSIZE),
			}
		}
		go goRoutine(str, myMap[str])
	}
	Debug("My map is ", myMap)
	f, _ := os.Open("aa")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ln := scanner.Text()

		for key, ch := range myMap {
			if strings.Contains(ln, key) {
				Gch = ch.mych
			}
		}
		processLine(ln, Gch)
	}
	if err := scanner.Err(); err != nil {
		Debug("Err:", err.Error())
	}
	stdscanner := bufio.NewScanner(os.Stdin)
	for {
		Debug("Enter 1 to do one")
		Debug("Enter 2 to do two ")
		Debug("Enter 3 to do threw")
		Debug("Enter 4 to do exit")
		stdscanner.Scan()
		ln := stdscanner.Text()
		Debug("Selected Choice is ", ln)
		if val, err := strconv.Atoi(ln); err == nil {
			if val == 4 {
				Debug("Exiting")
				return
			}
		} else {
			Debug("Error:", err.Error())
		}
	}
}

func processLine(ln string, mychan chan string) {
	mychan <- ln

}

func goRoutine(str string, mydata Data) {
	Debug("Starting Go Routine for str: ", str)
	for {
		select {
		case val := <-mydata.mych:
			Debug(str, ":", val)
			if mydata.funptr != nil {
				mydata.funptr(val, mydata)
			} else {
				Debug("NO processing defined for :", str)
			}
		}
	}
}
