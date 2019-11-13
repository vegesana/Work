package main

// TextProcessingTool.pdf
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var tempDir string
var myfile *os.File
var splfile *os.File

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

type pcldata struct {
	filename string
	line     string
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
type FInfo struct {
	modTime time.Time
	size    int64
}

var JunkCh chan interface{}

var FileDB map[string]FInfo

type GlobalData struct {
	// filename, vlan, mac, macInfo
	FDB      map[MacVlan]map[GPort]struct{}
	PortDB   map[GPort]map[MacVlan]struct{}
	VlanDB   map[GVlan]map[GPort]struct{}
	PinDB    map[GPort]GPort
	PclSlice []pcldata
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

func Error(param ...interface{}) {
	fmt.Fprintln(splfile, param)
	fmt.Println(param)
}

func Input(param ...interface{}) {
	fmt.Println(param)
}

func Dump(param ...interface{}) {
	//fmt.Fprintln(myfile, param)
	fmt.Println(param)
}

func Debug(param ...interface{}) {
	fmt.Fprintln(myfile, param)
	//fmt.Println(param)
}
func getFuncPtr(index int) func(lineData) {
	if index == 10 {
		return processMacInfo
	}
	if index == 8 {
		return processPclInfo
	}

	return nil
}
func main() {

	// String Slice - This will be adjusted based on the text in the
	// file
	tempDir = "/tmp/Raju/"
	defer myfile.Close()
	defer splfile.Close()

	myfile, _ = os.OpenFile(tempDir+"debug.log", os.O_WRONLY|os.O_CREATE, 0666)
	splfile, _ = os.OpenFile(tempDir+"error.log", os.O_WRONLY|os.O_CREATE, 0666)

	FileDB = map[string]FInfo{}
	Gdata = GlobalData{}
	Gdata.BufSize = 40
	Gdata.FDB = map[MacVlan]map[GPort]struct{}{}
	Gdata.PortDB = map[GPort]map[MacVlan]struct{}{}
	Gdata.VlanDB = map[GVlan]map[GPort]struct{}{}
	Gdata.PinDB = map[GPort]GPort{}
	Gdata.PclSlice = make([]pcldata, 0)
	Gdata.GWriteCh = make(chan interface{}, Gdata.BufSize)
	Gdata.GReadCh = make(chan interface{}) // Blocking Read

	go ProcessDataStructures()
	JunkCh = make(chan interface{}, Gdata.BufSize)

	go goRoutine("Junk", JunkCh, nil)

	// All the strings that are part of the input files. Each string
	// will get a channel created
	StrSlice = []string{"Network Interface Info", "NIC Interface Info",
		"NVMEOE Interface Info", "NIF Interface Info", "LIF Interface Info",
		"Statistics Info", "Cfg Info", "CPSS Info", "PCL Info",
		"VLAN Info", "MAC Info"}

	Gdata.MyMap = make(map[string]Data)
	for i, str := range StrSlice {
		tempCh := make(chan interface{}, Gdata.BufSize)
		fun := getFuncPtr(i)
		Gdata.MyMap[str] = Data{mych: tempCh, funptr: fun}
		go goRoutine(str, tempCh, fun)
	}

	Debug("My map is ", Gdata.MyMap)
	go func() {
		for {
			select {
			// This check whether anything got changed
			// in the file or any new file added ??
			case <-time.After(5 * time.Second):

				readAllFiles()
			}
		}
	}()
	readFromStdin()
}

func checkForLineDelimiter(ln string) bool {
	if strings.Contains(ln, "=======") {
		return true
	}
	return false
}
func isLineSpl(ln string) bool {
	if strings.Contains(ln, "CPSS Info") {
		return true
	}

	// Add all specical cases here
	return false

}
func readAllFiles() {
	// Host name is file name
	Debug("Reading all files")
	dir, _ := filepath.Glob(tempDir + "*.txt")
	for _, filename1 := range dir {
		// Check wether the file have changed - CRC/Checksum. If so
		// then run through it - ie First Remove the entries by that
		// filename
		filename := filepath.Base(filename1)
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))

		f, _ := os.Open(filename1)
		fi, _ := f.Stat()
		finfo := FInfo{fi.ModTime(), fi.Size()}

		if fvalue, ok := FileDB[filename]; ok {
			if finfo.modTime == fvalue.modTime &&
				finfo.size == fvalue.size {
				Debug("Same File - no action on file ", filename)
				continue
			}
		}

		// This create .tmp file in /same Dir for all txt files
		SplHandle()
		go splHandlingGo(tempDir)
		FileDB[filename] = finfo
		go func(f *os.File, filename string, mych chan interface{}) {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				ln := scanner.Text()
				for key, value := range Gdata.MyMap {
					// Special Write
					if strings.Contains(ln, key) {
						mych = value.mych
					}
				}
				lninfo := lineData{filename, ln, filename}
				mych <- lninfo
			}
			if err := scanner.Err(); err != nil {
				Debug("Err:", err.Error())
			}
		}(f, filename, JunkCh)
	}

}
func readFromStdin() {

	stdscanner := bufio.NewScanner(os.Stdin)
	for {
		Input("Enter 1 to dump macdb")
		Input("Enter 2 to dump pindb")
		Input("Enter 3 to dump pcldb")
		Input("Enter 4 to enter mac,vlan")
		Input("Enter 5 to servername,port")
		Input("Enter 7 to exit")
		stdscanner.Scan()
		ln := stdscanner.Text()
		Input("Selected Choice is ", ln)
		if val, err := strconv.Atoi(ln); err == nil {
			if val == 1 {
				Dump("Dump MAC db ")
				macvlan := MacVlan{}
				sendOnToReadCh(macvlan)
			}
			if val == 2 {
				Dump("Dump Pin db ")
				pin := PinInfo{}
				sendOnToReadCh(pin)
			}
			if val == 3 {
				Dump("Dump Pcl db ")
				pcl := pcldata{}
				sendOnToReadCh(pcl)
			}

			if val == 7 {
				Debug("Exiting")
				return
			}
		} else {
			Debug("Error:", err.Error())
		}
	}

}

func sendOnToReadCh(data interface{}) {
	Gdata.GReadCh <- data
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
	Dump("PrintfromDatabase")
	switch wdata.(type) {
	case MacVlan:
		Dump("MacVlan:Dump")
		for key, value := range Gdata.FDB {
			for key1, _ := range value {
				Dump(key, key1)
			}
		}
	case PinInfo:
		Dump("PinINfo:Dump")
		for key, value := range Gdata.PinDB {
			Dump(key, value)
		}
	case pcldata:
		Dump("Pcldata:Dump")
		for _, value := range Gdata.PclSlice {
			Dump(value)
		}

	}
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
		pindata := wdata.(PinInfo)
		Debug("Recevied Pindata", pindata)
		Gdata.PinDB[pindata.sport] = pindata.dport

	case pcldata:
		pata := wdata.(pcldata)
		Gdata.PclSlice = append(Gdata.PclSlice, pata)
	default:
		Debug("Writeing to database failed - Unknown type")
	}
}
func goRoutine(str string, ch chan interface{}, fun func(lineData)) {

	for {
		select {
		case val := <-ch:
			switch val.(type) {
			case lineData:
				myval := val.(lineData)
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