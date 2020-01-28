package main

// TextProcessingTool.pdf
import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tempDir string
var inputDir string
var myfile *os.File
var errfile *os.File

type SystemInfo struct {
	boardInfo string
	productId string
}
type MyError struct {
	myerr string
}

type MyInfo struct {
	myinfo string
}

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

var MyErrSlice []MyError
var MyInfoSlice []MyInfo
var JunkCh chan interface{}
var CounterCh chan interface{}
var NifInfoCh chan interface{}
var LifInfoCh chan interface{}
var MacInfoCh chan interface{}
var VlanInfoCh chan interface{}
var IntfInfoCh chan interface{}
var NvmInfoCh chan interface{}
var PclInfoCh chan interface{}
var StatsInfoCh chan interface{}
var CfgInfoCh chan interface{}
var NetworkInfoCh chan interface{}
var NicInfoCh chan interface{}

var FileDB map[string]FInfo
var TempIntfInfo map[string]string
var TempCntrInfo map[string]string
var TempStatsInfo map[string]string
var SystemMap map[string]SystemInfo

type GlobalData struct {
	// filename, vlan, mac, macInfo
	FDB         map[MacVlan]map[GPort]struct{}
	PortDB      map[GPort]map[MacVlan]struct{}
	VlanDB      map[GVlan]map[GPort]struct{}
	PinDB       map[GPort]GPort
	PclSlice    []pcldata
	BufSize     int
	MyMap       map[string]Data
	GWriteCh    chan interface{}
	GReadCh     chan interface{}
	ReadFilesCh chan interface{}
}

var Gdata GlobalData

type Data struct {
	mych   chan interface{}
	funptr func(lineData)
}

func Info(param ...interface{}) {
	fmt.Println(param)
}
func Error(param ...interface{}) {
	fmt.Fprintln(errfile, param)
	fmt.Println(param)
}
func Input(param ...interface{}) {
	fmt.Println(param)
}

func Dump(param ...interface{}) {
	fmt.Println(param)
}

func Debug(param ...interface{}) {
	fmt.Fprintln(myfile, param)
}

func resetDB() {
	Gdata.FDB = map[MacVlan]map[GPort]struct{}{}
	Gdata.PortDB = map[GPort]map[MacVlan]struct{}{}
	Gdata.VlanDB = map[GVlan]map[GPort]struct{}{}
	Gdata.PinDB = map[GPort]GPort{}
	MyErrSlice = nil
	MyInfoSlice = nil
}

func main() {

	tempDir = "/tmp/Raju/"
	inputDir = "/Users/rajuv/Techsupport/6052"
	SystemMap = map[string]SystemInfo{}
	defer myfile.Close()
	defer errfile.Close()

	pathValue := flag.String("path", "/Users/rajuv/Techsupport/6052",
		"Directoy Path")
	flag.Parse()
	fmt.Println("Curnt Path is ", *pathValue)
	inputDir = *pathValue

	myfile, _ = os.OpenFile(tempDir+"debug.log", os.O_WRONLY|os.O_CREATE, 0666)
	errfile, _ = os.OpenFile(tempDir+"error.log", os.O_WRONLY|os.O_CREATE, 0666)

	err := getFileContents("boardinfo.log")
	Debug("boardinfo.log file error is ", err)
	FileDB = map[string]FInfo{}
	TempIntfInfo = map[string]string{}
	TempCntrInfo = map[string]string{}
	TempStatsInfo = map[string]string{}

	Gdata = GlobalData{}
	Gdata.BufSize = 40
	Gdata.FDB = map[MacVlan]map[GPort]struct{}{}
	Gdata.PortDB = map[GPort]map[MacVlan]struct{}{}
	Gdata.VlanDB = map[GVlan]map[GPort]struct{}{}
	Gdata.PinDB = map[GPort]GPort{}
	Gdata.PclSlice = make([]pcldata, 0)
	Gdata.GWriteCh = make(chan interface{}, Gdata.BufSize)
	Gdata.GReadCh = make(chan interface{})     // Blocking Read
	Gdata.ReadFilesCh = make(chan interface{}) // Blocking Read

	go ProcessDataStructures()

	JunkCh = make(chan interface{}, Gdata.BufSize)
	CounterCh = make(chan interface{}, Gdata.BufSize)
	IntfInfoCh = make(chan interface{}, Gdata.BufSize)
	VlanInfoCh = make(chan interface{}, Gdata.BufSize)
	MacInfoCh = make(chan interface{}, Gdata.BufSize)
	NifInfoCh = make(chan interface{}, Gdata.BufSize)
	PclInfoCh = make(chan interface{}, Gdata.BufSize)
	StatsInfoCh = make(chan interface{}, Gdata.BufSize)
	CfgInfoCh = make(chan interface{}, Gdata.BufSize)
	NvmInfoCh = make(chan interface{}, Gdata.BufSize)
	LifInfoCh = make(chan interface{}, Gdata.BufSize)
	NetworkInfoCh = make(chan interface{}, Gdata.BufSize)
	NicInfoCh = make(chan interface{}, Gdata.BufSize)

	go goRoutine("Junk", JunkCh, nil)
	go goRoutine("Counters", CounterCh, CounterFun)
	go goRoutine("Interface Info", IntfInfoCh, IntfInfoFun)
	go goRoutine("Vlan Info", VlanInfoCh, VlanInfoFun)
	go goRoutine("Mac Info", MacInfoCh, MacInfoFun)
	go goRoutine("Nif Info", NifInfoCh, NifInfoFun)
	go goRoutine("Pcl Info", PclInfoCh, PclInfoFun)
	go goRoutine("Stats Info", StatsInfoCh, StatsInfoFun)
	go goRoutine("Cfg Info", CfgInfoCh, CfgInfoFun)
	go goRoutine("Nvm Info", NvmInfoCh, NvmInfoFun)
	go goRoutine("Lif Info", LifInfoCh, LifInfoFun)
	go goRoutine("Netowrk Info", NetworkInfoCh, NetworkInfoFun)
	go goRoutine("Nic Info", NicInfoCh, NicInfoFun)

	go func() {
		Debug("Go routine to read ncdutil file")
		for {
			select {
			// This check whether anything got changed
			// in the file or any new file added ??
			case <-Gdata.ReadFilesCh:
				readLogFiles("ncdutil.log")
			}
		}
	}()
	for server, sysinfo := range SystemMap {
		fmt.Println(server, sysinfo)
		if strings.Contains(sysinfo.boardInfo, "Boston") {
			Gdata.ReadFilesCh <- struct{}{}
		} else {
			errstr := fmt.Sprintf("%s is NOT D20\n", server)
			SendError(errstr)
		}
	}
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
func getServerNameFromPath(path string) string {
	re, _ := regexp.Compile(`appserv\d+`)
	value := re.FindString(path)
	return value
}

func getFileContents(filename string) error {

	var servername string
	err := errors.New("File not found")

	tempDir := map[string]struct{}{}
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {

		if info.Name() == filename {
			err = nil
			if servername = getServerNameFromPath(path); servername == "" {
				Error("No Name with appserv found in path", path)
				return nil
			}
			if _, ok := tempDir[servername]; ok {
				fmt.Println("File path already exits", path)
				Debug("File path already exits", path)
				return nil
			}
			tempDir[servername] = struct{}{}
			f, _ := os.Open(path)
			allText, _ := ioutil.ReadAll(f)
			productid := getValueOfStr(string(allText), "Product id", ":")
			board := getValueOfStr(string(allText), "Board", ":")

			SystemMap[servername] = SystemInfo{board, productid}

			return nil
		}
		return nil
	})
	return err
}
func readLogFiles(filename string) error {
	var servername string
	// Host name is file name
	Debug("Reading all NCD files")
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Error("FIlepaqqth walk failed")
			return err
		}

		if info.Name() == filename {
			if servername = getServerNameFromPath(path); servername == "" {
				Error("No Name with appserv found in path", path)
				return nil
			}
			Debug("File Path is ", path)
			f, _ := os.Open(path)
			go func(f *os.File, servername string, mych chan interface{}) {
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					origTxt := scanner.Text()
					ln := strings.TrimSpace(origTxt)

					if strings.Index(ln, "Counters Info:") == 0 {
						Debug("Picked Counter Channel")
						mych = CounterCh
					}

					if strings.Index(ln, "QOS Info:") == 0 {
						Debug("Picked Qos Channel")
						mych = CounterCh
					}

					if strings.Index(ln, "Network Interface Info") == 0 {
						Debug("Picked Network Interface Channel")
						mych = NetworkInfoCh
					}

					if strings.Index(ln, "NIC Interface Info") == 0 {
						Debug("Picked Network Interface Channel")
						mych = NicInfoCh
					}

					if strings.Index(ln, "NVMEOE Interface Info") == 0 {
						Debug("NVMEOE Interface Info channel ")
						mych = NvmInfoCh
					}
					if strings.Index(ln, "VLAN Info") == 0 {
						Debug("Vlan Info channel ")
						mych = VlanInfoCh
					}
					if strings.Index(ln, "MAC Info") == 0 {
						Debug("MAC Info channel ")
						mych = MacInfoCh
					}
					if strings.Index(ln, "PCL Info") == 0 {
						Debug("PCL Info channel ")
						mych = PclInfoCh
					}
					if strings.Index(ln, "Interface Info") == 0 {
						Debug("Interface Info channel ")
						mych = IntfInfoCh
					}
					if strings.Index(ln, "LIF Interface Info") == 0 {
						Debug("Lif interface Info channel ")
						mych = LifInfoCh
					}
					if strings.Index(ln, "NIF Interface Info") == 0 {
						Debug("Lif interface Info channel ")
						mych = NifInfoCh
					}

					if strings.Index(ln, "Statistics Info") == 0 {
						Debug("Statistic Info channel ")
						mych = StatsInfoCh
					}
					if strings.Index(ln, "Cfg Info") == 0 {
						Debug("Cfg Info channel ")
						mych = CfgInfoCh
					}

					lninfo := lineData{servername, origTxt, servername}
					mych <- lninfo

					// All the end to determine the end of Channel
					if strings.Contains(ln, "maxPortDescrLimit") {
						// Place HOlder till we find new channel
						mych = JunkCh
					}

				}
				if err := scanner.Err(); err != nil {
					Debug("Err:", err.Error())
				}
			}(f, servername, JunkCh)
		}
		return nil
	})
	return nil
}
func readFromStdin() {

	// Keyboard Input
	stdscanner := bufio.NewScanner(os.Stdin)
	for {
		time.Sleep(time.Second * 2)
		Input("Enter 1 to dump macdb")
		Input("Enter 2 to dump pindb")
		Input("Enter 3 to dump pcldb")
		Input("Enter 4 to Dump all Errors")
		Input("Enter 5 to Read All input files again")
		Input("Enter 6 to  Dump System Platform Info")
		Input("Enter 7 to exit")
		stdscanner.Scan()
		ln := stdscanner.Text()
		Input("Selected Choice is ", ln)
		if strings.TrimSpace(ln) == "" {
			continue
		}
		if val, err := strconv.Atoi(ln); err == nil {

			switch val {
			case 1:
				Dump("Dump MAC db ")
				macvlan := MacVlan{}
				sendOnToReadCh(macvlan)
			case 2:
				Dump("Dump Pin db ")
				pin := PinInfo{}
				sendOnToReadCh(pin)
			case 3:
				Dump("Dump Pcl db ")
				pcl := pcldata{}
				sendOnToReadCh(pcl)
			case 4:
				Dump("Dump Errors from File")
				DumpErrors()
			case 5:
				Dump("Read all files again ..Reset all DS")
				resetDB()
				Gdata.ReadFilesCh <- struct{}{}
			case 6:
				Dump("System Info Dump")
				sysinfo := SystemInfo{}
				sendOnToReadCh(sysinfo)

			case 7:
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
	case SystemInfo:
		Dump("SystemInfo:Dump")
		for key, value := range SystemMap {
			Dump(key, value)
		}

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
	case MyError:
		errdata := wdata.(MyError)
		MyErrSlice = append(MyErrSlice, errdata)
	case MyInfo:
		infodata := wdata.(MyInfo)
		MyInfoSlice = append(MyInfoSlice, infodata)

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
func (line lineData) getText() string {
	return line.line
}
func (line lineData) getFileName() string {
	return line.filename
}
