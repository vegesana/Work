package main

// TextProcessingTool.pdf
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var myfile *os.File

type MacInfo struct {
	fileName string
	mac      string
	vlan     string
	isStatic string
	portNum  string
	Index    int
}
type lineData struct {
	filename string
	line     string
	hostname string
}

// Map index
type MacVlan struct {
	mac  string
	vlan string
}
type GPort struct {
	filename string
	portname string
}
type PinInfo struct {
	sport GPort
	dport GPort
}

type GVlan struct {
	filename string
	vlan     string
}

var JunkCh chan interface{}

type GlobalData struct {
	// filename, vlan, mac, macInfo
	FDB      map[MacVlan]map[GPort]struct{}
	PortDB   map[GPort]map[MacVlan]struct{}
	VlanDB   map[GVlan]map[GPort]struct{}
	PinDB    map[GPort]GPort
	BufSize  int
	MyMap    map[string]Data
	GWriteCh chan interface{}
	GReadCh  chan interface{}
}

var Gdata GlobalData

type Data struct {
	mych   chan interface{}
	funptr func(lineData)
}

var StrSlice []string

func Debug(param ...interface{}) {
	fmt.Fprintln(myfile, param)
	fmt.Println(param)
}
func getFuncPtr(index int) func(lineData) {
	if index == 10 {
		return processMacInfo
	}
	if index == 7 { // overloading CPSS info for Interface Info
		return processPinInfo
	}
	return nil
}
func main() {

	// String Slice - This will be adjusted based on the text in the
	// file
	defer myfile.Close()

	myfile, _ = os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE, 0666)

	Gdata = GlobalData{}
	Gdata.BufSize = 40
	Gdata.FDB = map[MacVlan]map[GPort]struct{}{}
	Gdata.PortDB = map[GPort]map[MacVlan]struct{}{}
	Gdata.VlanDB = map[GVlan]map[GPort]struct{}{}
	Gdata.PinDB = map[GPort]GPort{}
	Gdata.GWriteCh = make(chan interface{}, Gdata.BufSize)
	Gdata.GReadCh = make(chan interface{}) // Blocking Read

	go ProcessDataStructures()
	JunkCh = make(chan interface{}, Gdata.BufSize)

	go goRoutine("Junk", JunkCh)

	// All the strings that are part of the input files. Each string
	// will get a channel created
	StrSlice = []string{"Network Interface Info", "NIC Interface Info",
		"NVMEOE Interface Info", "NIF Interface Info", "LIF Interface Info",
		"Statistics Info", "Cfg Info", "CPSS Info", "PCL Info",
		"VLAN Info", "MAC Info"}

	Gdata.MyMap = make(map[string]Data)
	for i, str := range StrSlice {
		tempCh := make(chan interface{}, Gdata.BufSize)
		Gdata.MyMap[str] = Data{mych: tempCh, funptr: getFuncPtr(i)}
	}
	for _, str := range StrSlice {
		go goRoutine(str, Gdata.MyMap[str].mych)
	}

	Debug("My map is ", Gdata.MyMap)
	readAllFiles()
	readFromStdin()
}

func readAllFiles() {
	// Host name is file name
	dir, _ := filepath.Glob("*.txt")
	for _, name := range dir {
		// Check wether the file have changed - CRC/Checksum. If so
		// then run through it - ie First Remove the entries by that
		// filename
		go func(filename string, mych chan interface{}) {
			f, _ := os.Open(filename)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				ln := scanner.Text()

				for key, value := range Gdata.MyMap {
					if strings.Contains(ln, key) {
						mydata := value
						mych = mydata.mych
					}
				}
				lninfo := lineData{filename, ln, filename}
				mych <- lninfo
			}
			if err := scanner.Err(); err != nil {
				Debug("Err:", err.Error())
			}
		}(name, JunkCh)
	}

}
func readFromStdin() {

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

func ProcessDataStructures() {
	for {
		select {
		case wdata := <-Gdata.GWriteCh:
			buildDatabase(wdata)
		case rdata := <-Gdata.GReadCh:
			printFromDatabase(rdata)
		}
	}

}

func printFromDatabase(wdata interface{}) {
}

func buildDatabase(wdata interface{}) {
	switch wdata.(type) {
	case MacInfo:
		macdata := wdata.(MacInfo)
		macvlankey := MacVlan{mac: macdata.mac, vlan: macdata.vlan}
		subkey := GPort{macdata.fileName, macdata.portNum}
		if value, ok := Gdata.FDB[macvlankey]; !ok {
			Gdata.FDB[macvlankey] = make(map[GPort]struct{})
			Gdata.FDB[macvlankey][subkey] = struct{}{}
		} else {
			if newval, ok := value[subkey]; !ok {
				value[subkey] = struct{}{}
				Gdata.FDB[macvlankey] = value
			} else {
				Debug("Value already part of map", newval)
			}
		}
		Debug("BUILD MACINFO", Gdata.FDB)

		if value, ok := Gdata.PortDB[subkey]; !ok {
			Gdata.PortDB[subkey] = make(map[MacVlan]struct{})
			Gdata.PortDB[subkey][macvlankey] = struct{}{}
		} else {
			if newval, ok := value[macvlankey]; !ok {
				value[macvlankey] = struct{}{}
				Gdata.PortDB[subkey] = value
			} else {
				Debug("Value already part of map", newval)
			}
		}
		Debug("BUILD PortDB", Gdata.PortDB)
	case PinInfo:

	default:
		Debug("Writeing to database failed - Unknown type")
	}
}
func goRoutine(str string, ch chan interface{}) {
	var fun func(lineData)

	if val, ok := Gdata.MyMap[str]; ok {
		fun = val.funptr
	}
	for {
		select {
		case val := <-ch:
			switch val.(type) {
			case lineData:
				myval := val.(lineData)
				Debug("line:", myval.line)
				Debug("filename:", myval.filename)
				Debug("hostname:", myval.hostname)
				// Check wheather we have any function pointer
				// for the given
				if fun != nil {
					fun(myval)
				} else {
					Debug("NO processing defined for :")
				}

			default:
			}

		}
	}
}
