package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	Server   string
	Mac      string
	Vlan     string
	IsStatic string
	PortNum  string // THis include servername and Port Number
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
	Server   string
	Line     string
	hostname string
}
type Pcldata struct {
	Server string
	Line   string
}

// Map index
type MacVlan struct {
	Mac  string
	Vlan string
}
type GPort struct {
	Server   string
	Portname string
}
type PinInfo struct {
	Server string
	Sport  string
	Dport  string
}

type GVlan struct {
	Server string
	vlan   string
}
type FInfo struct {
	modTime time.Time
	size    int64
}
type GlobalData struct {
	FDB       map[MacVlan]map[GPort]struct{}
	PortDB    map[GPort]map[MacVlan]struct{}
	VlanDB    map[GVlan]map[GPort]struct{}
	NodeDB    NodeList
	NetworkDB NetworkList
	// appser33: map {1-3 4,5}
	PinDB       map[string]map[string]string
	PclDb       map[string][]string
	BufSize     int
	ReadFilesCh chan interface{}
}

var Gdata GlobalData
var btpUtilMap map[Key]string
var AllStringKeyMap map[string]struct{}

var JunkCh chan interface{}
var NetworkInfoCh chan interface{}
var MacInfoCh chan interface{}
var IntfInfoCh chan interface{}
var PclInfoCh chan interface{}
var CounterCh chan interface{}
var StatsInfoCh chan interface{}
var CfgInfoCh chan interface{}

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
var TempIntfInfo map[string]string
var TempCntrInfo map[string]string
var TempStatsInfo map[string]string

const (
	SIZE = 10
)

var MyErrSlice []MyError
var tempErrorExistsMap map[MyError]struct{}

func init() {
	fmt.Println("Debug Server Init called")
	SystemMap = map[string]SystemInfo{}
	TempIntfInfo = map[string]string{}
	TempCntrInfo = map[string]string{}
	TempStatsInfo = map[string]string{}
	tempErrorExistsMap = map[MyError]struct{}{}
	dbWriteCh = make(chan interface{})
	//dbReadCh = make(chan ReadData, SIZE)
	dbReadCh = make(chan ReadData, SIZE)
	MyErrSlice = make([]MyError, 0)

	Gdata = GlobalData{}
	Gdata.BufSize = 20

	Gdata.FDB = map[MacVlan]map[GPort]struct{}{}
	Gdata.PortDB = map[GPort]map[MacVlan]struct{}{}
	Gdata.VlanDB = map[GVlan]map[GPort]struct{}{}
	Gdata.PinDB = map[string]map[string]string{}
	Gdata.PclDb = map[string][]string{}

	JunkCh = make(chan interface{}, Gdata.BufSize)
	NetworkInfoCh = make(chan interface{}, Gdata.BufSize)
	MacInfoCh = make(chan interface{}, Gdata.BufSize)
	IntfInfoCh = make(chan interface{}, Gdata.BufSize)
	PclInfoCh = make(chan interface{}, Gdata.BufSize)
	CounterCh = make(chan interface{}, Gdata.BufSize)
	StatsInfoCh = make(chan interface{}, Gdata.BufSize)
	CfgInfoCh = make(chan interface{}, Gdata.BufSize)

	go readWriteGoRoutine()

	go goRoutine("Junk", JunkCh, nil)
	go goRoutine("Network Info", NetworkInfoCh, NetworkInfoFun)
	go goRoutine("Mac Info", MacInfoCh, MacInfoFun)
	go goRoutine("Interface Info", IntfInfoCh, IntfInfoFun)
	go goRoutine("Pcl Info", PclInfoCh, PclInfoFun)
	go goRoutine("Counters", CounterCh, CounterFun)
	go goRoutine("Stats Info", StatsInfoCh, StatsInfoFun)
	go goRoutine("Cfg Info", CfgInfoCh, CfgInfoFun)

	// Start a Go routine to update the DATABASES and also
	// to get the data from DB
}
func Start(path string, serverprefix string) {
	fmt.Println("Debug Server")
	go loopThroughAllFilesInAllSubDir(path, serverprefix)
}

func getServerNameFromPath(path string, servprefix string) string {
	newstr := servprefix + `\d+`
	re, _ := regexp.Compile(newstr)
	value := re.FindString(path)
	return value
}

func loopThroughAllFilesInAllSubDir(inputDir string, servprefix string) error {

	var servername string
	var oneTimeNetworkCfg bool

	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {

		if strings.Contains(path, servprefix) {
			if servername = getServerNameFromPath(path, servprefix); servername != "" {

				/* If boardinfo exits, then process it */
				if info.Name() == "boardinfo.log" {
					fmt.Println("Boardinfo ", path)
					processBoardInfo(path, servername)
				}
				if info.Name() == "fwdcounters.log" {
					fmt.Println("fwdcounter", path)
					processFwdCounters(path, servername)
				}
				if info.Name() == "btputil.log" {
					fmt.Println("btputil", path)
					processBtpUtil(path, servername)
				}

				if info.Name() == "ncdutil.log" {
					go processNcdUtil(path, servername)
				}

				if info.Name() == "node.log" {
					go ProcessNodeInfo(path, servername)
				}
				if info.Name() == "network.log" {
					if !oneTimeNetworkCfg {
						go ProcessNetworkInfo(path, servername)
						oneTimeNetworkCfg = true
					}
				}

			}
		}
		return nil
	})
	return nil
}

func checkErrorExists(myerr MyError) bool {

	if _, ok := tempErrorExistsMap[myerr]; !ok {
		tempErrorExistsMap[myerr] = struct{}{}
		return false
	}
	return true
}
func writeToDBBackend(wval interface{}) {
	switch wval.(type) {
	case NodeList:
		nval := wval.(NodeList)
		Gdata.NodeDB = nval
	case NetworkList:
		nval := wval.(NetworkList)
		Gdata.NetworkDB = nval
	case SystemInfo:
		wsval := wval.(SystemInfo)
		fmt.Printf("WriteToDB System Info: %#v\n", wsval)
		if len(wsval.ServerName) == 0 {
			fmt.Println("Serever INFO NOT FOUND ", SystemMap)
		} else {
			SystemMap[wsval.ServerName] = SystemInfo{wsval.BoardInfo,
				wsval.ProductId, wsval.ServerName}
		}
	case MyError:
		myerr := wval.(MyError)

		if !checkErrorExists(myerr) {
			fmt.Printf("WritetToDB MyError: %#v\n", myerr)
			MyErrSlice = append(MyErrSlice, myerr)
		}
	case MacInfo:
		/* we are buidling 2 data based here */
		macdata := wval.(MacInfo)
		macvlankey := MacVlan{Mac: macdata.Mac, Vlan: macdata.Vlan}
		subkey := GPort{macdata.Server, macdata.PortNum}
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
		pindata := wval.(PinInfo)
		if value, ok := Gdata.PinDB[pindata.Server]; !ok {
			Gdata.PinDB[pindata.Server] = make(map[string]string)
			Gdata.PinDB[pindata.Server][pindata.Sport] = pindata.Dport
		} else {
			if value[pindata.Dport] != pindata.Sport {
				value[pindata.Sport] = pindata.Dport
				Gdata.PinDB[pindata.Server] = value
			}
		}
	case Pcldata:
		pata := wval.(Pcldata)

		Gdata.PclDb[pata.Server] = append(Gdata.PclDb[pata.Server],
			pata.Line)
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
	case MacInfo:
		/* We can make a copy of it and send */
		outIntf = Gdata.FDB
	case PinInfo:
		outIntf = Gdata.PinDB
	case Pcldata:
		outIntf = Gdata.PclDb
	case NetworkList:
		outIntf = Gdata.NetworkDB
	case NodeList:
		outIntf = Gdata.NodeDB
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
func GetPclInfo() interface{} {
	fmt.Println("Geting Pcl informatin")
	pclinfo := Pcldata{}
	return readFromDb(pclinfo)
}
func GetPinInfo() interface{} {
	fmt.Println("Geting Pin informatin")
	pininfo := PinInfo{}
	return readFromDb(pininfo)
}

func GetMacInfo() interface{} {
	fmt.Println("Geting Mac informatin")
	macinfo := MacInfo{}
	return readFromDb(macinfo)
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
