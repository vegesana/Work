package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Key struct {
	lport, session, key string
}

// valid session
type VSession struct {
	lport   string
	session string
}

type MacInfo struct {
	fileName string
	mac      string
	vlan     string
	isStatic string
	portNum  string
	Index    int
}
type SystemInfo struct {
	BoardInfo  string
	ProductId  string
	ServerName string
}
type ReadData struct {
	data  interface{}
	rchan chan interface{}
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
type GlobalData struct {
	// filename, vlan, mac, macInfo
	FDB         map[MacVlan]map[GPort]struct{}
	PortDB      map[GPort]map[MacVlan]struct{}
	VlanDB      map[GVlan]map[GPort]struct{}
	PinDB       map[GPort]GPort
	PclSlice    []pcldata
	BufSize     int
	ReadFilesCh chan interface{}
}

var Gdata GlobalData
var btpUtilMap map[Key]string
var AllStringKeyMap map[string]struct{}

var JunkCh chan interface{}
var NetworkInfoCh chan interface{}
var MacInfoCh chan interface{}

type MyError struct {
	ServerName string
	MyErr      string
}
type MyInfo struct {
	myinfo string
}

var dbWriteCh chan interface{}
var dbReadCh chan ReadData

var SystemMap map[string]SystemInfo

const (
	SIZE = 10
)

var MyErrSlice []MyError

func init() {
	fmt.Println("Debug Server Init called")
	SystemMap = map[string]SystemInfo{}
	dbWriteCh = make(chan interface{})
	dbReadCh = make(chan ReadData, SIZE)
	MyErrSlice = make([]MyError, 0)

	Gdata = GlobalData{}
	Gdata.BufSize = 20

	Gdata.FDB = map[MacVlan]map[GPort]struct{}{}
	Gdata.PortDB = map[GPort]map[MacVlan]struct{}{}
	Gdata.VlanDB = map[GVlan]map[GPort]struct{}{}

	JunkCh = make(chan interface{}, Gdata.BufSize)
	NetworkInfoCh = make(chan interface{}, Gdata.BufSize)
	MacInfoCh = make(chan interface{}, Gdata.BufSize)

	go readWriteGoRoutine()

	go goRoutine("Junk", JunkCh, nil)
	go goRoutine("Network Info", NetworkInfoCh, NetworkInfoFun)
	go goRoutine("Mac Info", MacInfoCh, MacInfoFun)

	// Start a Go routine to update the DATABASES and also
	// to get the data from DB
}
func Start(path string) {
	fmt.Println("Debug Server")
	go loopThroughAllFilesInAllSubDir(path)
}

func getServerNameFromPath(path string) string {
	re, _ := regexp.Compile(`appserv\d+`)
	value := re.FindString(path)
	return value
}

func loopThroughAllFilesInAllSubDir(inputDir string) error {

	var servername string

	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {

		if servername = getServerNameFromPath(path); servername != "" {

			/* If boardinfo exits, then process it */
			if info.Name() == "boardinfo.log" {
				go processBoardInfo(path, servername)
			}
			if info.Name() == "fwdcounters.log" {
				go processFwdCounters(path, servername)
			}

			if info.Name() == "btputil.log" {
				go processBtpUtil(path, servername)
			}
			/*
				if info.Name() == "ncdutil.log" {
					go processNcdUtil(path, servername)
				}
			*/
		}
		return nil
	})
	return nil
}

func writeToDBBackend(wval interface{}) {
	switch wval.(type) {
	case SystemInfo:
		wsval := wval.(SystemInfo)
		fmt.Println("Write System Info", SystemMap)
		SystemMap[wsval.ServerName] = SystemInfo{wsval.BoardInfo,
			wsval.ProductId, wsval.ServerName}
	case MyError:
		myerr := wval.(MyError)
		MyErrSlice = append(MyErrSlice, myerr)
	case MacInfo:
		/* we are buidling 2 data based here */
		macdata := wval.(MacInfo)
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
		fmt.Println("Bulding MAC and FDB DB DONE")
	default:
		fmt.Println("Write of unknown type ")
	}

}

func readFromDb(val interface{}) interface{} {
	routputch := make(chan interface{})
	rdata := ReadData{val, routputch}
	dbReadCh <- rdata
	return <-routputch
}

func writeToDb(val interface{}) {
	dbWriteCh <- val
}

func workOnReadFromDBChan(rval ReadData) {

	var outIntf interface{}
	data := rval.data
	outch := rval.rchan

	defer func() {
		outch <- outIntf
	}()

	switch data.(type) {
	case SystemInfo:
		copymap := map[string]SystemInfo{}
		for key, value := range SystemMap {
			copymap[key] = value
		}
		outIntf = copymap
	case MyError:
		newslice := MyErrSlice
		outIntf = newslice
	default:
		fmt.Println("Read for unknown type of data")
	}
}
func readWriteGoRoutine() {
	for {
		select {
		case wval := <-dbWriteCh:
			writeToDBBackend(wval)
		case rval := <-dbReadCh:
			workOnReadFromDBChan(rval)
		}

	}
}

func GetMacInfo() interface{} {
	fmt.Println("Geting Mac informatin")
	sysinfo := SystemInfo{}
	return readFromDb(sysinfo)
}

// this should send read on to the channel to get data
func GetErrorInfo() interface{} {
	fmt.Println("Geting Error informatin")
	err := MyError{}
	return readFromDb(err)
}

// this should send read on to the channel to get data
func GetSystemInfo() interface{} {
	fmt.Println("Geting system informatin")
	sysinfo := SystemInfo{}
	return readFromDb(sysinfo)
}
