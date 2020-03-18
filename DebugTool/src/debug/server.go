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

type SystemInfo struct {
	BoardInfo  string
	ProductId  string
	ServerName string
}
type ReadData struct {
	data  interface{}
	rchan chan interface{}
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

	go readWriteGoRoutine()

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

			if info.Name() == "ncdutil.log" {
				processNcdUtil(path, servername)
			}
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
