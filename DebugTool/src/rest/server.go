package rest

import (
	"DebugTool/src/debug"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type page struct {
	Title   string
	SysInfo []debug.SystemInfo
	MacInfo []MacInfoRest
	PinInfo []PinInfoRest
	PclInfo []PclInfoRest
	ErrInfo []debug.MyError
	// This can be Generic - Copy and paste
	NetworkInfor []NetworkRest
	NodeInfor    []NodeRest
	// Map[string][]map[string]string i.e servername: values id:key,value
	CtrlInfor    []CtrlInfoRest
	D10lifSlice  []D10LifInfoRest
	D10VlanSlice []D10VlanInfoRest
}
type CtrlInfoRest struct {
	Servername   string
	Id           string
	Type         string
	LocalMac     string
	RemoteMac    string
	PairedCtrlId string
	Cookie       string
}
type NetworkRest struct {
	Name       string // Name
	Zone       string // Annotations[failure-domain.beta.kubernetes.io/zone]
	Type       string // Spec.Type
	SSubnet    string // Status.Subnet
	SVlan      string // Status.Vlan [int to string]
	SUsedAddr  string // Status.UsedAddrs
	SNumAddr   string // Status.NumAddrs
	GatewayMac string // StorageNetworkSpec.GatewayMac: only when stroage
}

type NodeRest struct {
	HostName   string // Host.Hostname
	NodeHealth string // Status.NodeHealth
	K8sHealth  string // Status.K8sHealth
	Zone       string // Spec.Zone
	SIP1       string // StorageIPs.StorageIpPort1
	SIP3       string // StorageIPs.StorageIpPort2
	GMAC       string // StorageIPs.GatewayMac
	Svlan      string // StorageIPs.SVLAN
	Mode       string // StorageIPs.Mode
}
type RuleListRest struct {
	Rule string
}
type PclInfoRest struct {
	Server   string
	RuleList []RuleListRest
}

type PortPairRest struct {
	SPort string
	DPort string
}
type D10VlanRest struct {
	Vlan      uint64
	LifBitmap string
}
type D10VlanInfoRest struct {
	Server  string
	D10Vlan []D10VlanRest
}

type LifRest struct {
	LifId uint64
	Vlan  uint64
	Pif   uint64
}
type D10LifInfoRest struct {
	Server   string
	LifSlice []LifRest
}

type PinInfoRest struct {
	Server    string
	PortSlice []PortPairRest
}
type ServerPortInfo struct {
	Server string
	Port   string
}
type MacInfoRest struct {
	Mac             string
	Vlan            string
	ServerPortSlice []ServerPortInfo
}
type RestObj struct {
}

func init() {
	fmt.Println("Default Packet called API ")
}

func Init() *RestObj {
	return &RestObj{}
}
func (r *RestObj) Start() {
	fmt.Println("Rest Start")
	/* use default Mux (nil)  */
	http.HandleFunc("/", r.HandleMainConfig)
	http.HandleFunc("/debugsubmit", r.HandleDebugConfig)
	http.HandleFunc("/testcasesubmit", r.HandleTestcases)
	http.HandleFunc("/sysSubmit", r.HandleSystemInfo)
	http.HandleFunc("/macSubmit", r.HandleMacInfo)
	http.HandleFunc("/errSubmit", r.HandleErrInfo)
	http.HandleFunc("/pclSubmit", r.HandlePclInfo)
	http.HandleFunc("/pinSubmit", r.HandlePinInfo)
	http.HandleFunc("/networkSubmit", r.HandleNetworkInfo)
	http.HandleFunc("/nodeSubmit", r.HandleNodeInfo)
	http.HandleFunc("/sputilSubmit", r.HandleSputilInfo)
	http.HandleFunc("/ctrlSubmit", r.HandleCtrlInfo)
	go http.ListenAndServe("Localhost:8080", nil)
}

/*
	The HandleMainConfig handler will generte HTML file tht goes to the
	client (brower). Once user clicks on thte button, we write our HTML
	in such a ways that it calls specific webservice (/debugsubmit)
	form action="/debugsubmit
*/
func (r *RestObj) HandleMainConfig(resp http.ResponseWriter, req *http.Request) {
	filepath := "config/main.html"
	http.ServeFile(resp, req, filepath)
}
func (r *RestObj) HandleNodeInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetNodeInfo()
	nodeinfo := val.(debug.NodeList)

	w.Header().Add("Content Type", "text/html")
	templates := template.New("template")
	templates.New("Body").Parse(nodeDoc)
	templates.New("List").Parse(nodeDocList)

	nodeslice := []NodeRest{}

	for _, nl := range nodeinfo.Items {
		nodeelem := NodeRest{} // Fill the data

		nodeelem.HostName = nl.Host.Hostname
		nodeelem.NodeHealth = string(nl.Status.NodeHealth)
		nodeelem.K8sHealth = string(nl.Status.K8sHealth)
		nodeelem.Zone = nl.Spec.Zone
		nodeelem.SIP1 = nl.StorageIPs.StorageIpPort1
		nodeelem.SIP3 = nl.StorageIPs.StorageIpPort3
		nodeelem.GMAC = nl.StorageIPs.GatewayMac
		nodeelem.Svlan = strconv.Itoa(int(nl.StorageIPs.SVLAN))
		nodeelem.Mode = string(nl.StorageIPs.Mode)
		nodeslice = append(nodeslice, nodeelem)
	}

	fmt.Printf("NodeInfo %#v\n", nodeslice)
	page := page{Title: "Node Information", NodeInfor: nodeslice}
	templates.Lookup("Body").Execute(w, page)
	return
}

func (r *RestObj) HandleSputilInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetSputilInfo()
	fmt.Printf("HandleSputilInfo", val)
}

// For any []{servername, ,map[int]mapp[string]string)
func (r *RestObj) HandleCtrlInfo(w http.ResponseWriter, req *http.Request) {
	info := debug.GetCtrlInfo()

	infoSlice := info.([]debug.NameMap)

	w.Header().Add("Content Type", "text/html")
	templates := template.New("template")
	templates.New("Body").Parse(ctrlDoc)
	templates.New("List").Parse(ctrlDocList)

	ctrlslice := []CtrlInfoRest{}

	// Flatenning here : New
	for _, cinfo := range infoSlice {

		for _, valuemap := range cinfo.Mymap {
			/* Right elwement from backend dump info*/
			ctrlelem := CtrlInfoRest{} // Fill the data
			ctrlelem.Servername = cinfo.Servername
			ctrlelem.Id = valuemap["id"]
			var mytype string
			switch strings.TrimSpace(valuemap["type"]) {
			case "2":
				mytype = "PROXY"
			case "4":
				mytype = "REMOTE"
			case "1":
				mytype = "LOCAL"
			default:
				mytype = "UNKNOWN"
			}
			ctrlelem.Type = mytype
			ctrlelem.LocalMac = valuemap["local_mac"]
			ctrlelem.RemoteMac = valuemap["remote_mac"]
			ctrlelem.PairedCtrlId = valuemap["paired_ctrl_id"]
			ctrlelem.Cookie = valuemap["ctrl shared cookie"]
			ctrlslice = append(ctrlslice, ctrlelem)
		}

	}
	fmt.Printf("CtrlInfo %#v\n", ctrlslice)
	page := page{Title: "Control Information", CtrlInfor: ctrlslice}
	templates.Lookup("Body").Execute(w, page)

}

func (r *RestObj) HandleNetworkInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetNetworkInfo()

	ninfo := val.(debug.NetworkList)
	w.Header().Add("Content Type", "text/html")
	templates := template.New("template")
	templates.New("Body").Parse(networkDoc)
	templates.New("List").Parse(networkDocList)

	networkslice := []NetworkRest{}

	for _, nl := range ninfo.Items {
		networkelem := NetworkRest{} // Fill the data
		networkelem.Name = nl.Name
		networkelem.Zone = nl.Annotations["failure-domain.beta.kubernetes.io/zone"]
		networkelem.Type = nl.Spec.Type
		networkelem.SSubnet = nl.Status.Subnet
		networkelem.SVlan = strconv.Itoa(int(nl.Status.VLAN))
		networkelem.SUsedAddr = strconv.Itoa(nl.Status.UsedAddrs)
		networkelem.SNumAddr = strconv.Itoa(nl.Status.NumAddrs)
		networkelem.GatewayMac = nl.StorageNetworkSpec.GatewayMac

		networkslice = append(networkslice, networkelem)
	}

	fmt.Printf("NetworkInfo %#v\n", networkslice)
	page := page{Title: "Network Information", NetworkInfor: networkslice}
	templates.Lookup("Body").Execute(w, page)

	return
}

// Map[sgring]map[string]string
func (r *RestObj) HandlePinInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetPinInfo()
	pininfo := val.(map[string]map[string]string)
	fmt.Printf("pin info rest%#v\n", pininfo)

	w.Header().Add("Content Type", "text/html")
	templates := template.New("template")
	templates.New("Body").Parse(pinDoc)
	templates.New("List").Parse(pinDocList)
	templates.New("List1").Parse(pinDocList1)

	pinInfoRest := []PinInfoRest{}
	for server, portMap := range pininfo {
		pininfo := PinInfoRest{}
		pininfo.Server = server
		for sport, dport := range portMap {
			pininfoport := PortPairRest{}
			pininfoport.SPort = sport
			pininfoport.DPort = dport
			pininfo.PortSlice = append(pininfo.PortSlice, pininfoport)
		}
		pinInfoRest = append(pinInfoRest, pininfo)
	}

	page := page{Title: "Pin Information", PinInfo: pinInfoRest}
	templates.Lookup("Body").Execute(w, page)

}

func (r *RestObj) HandlePclInfo(w http.ResponseWriter, req *http.Request) {

	val := debug.GetPclInfo()
	pclInfo := val.(map[string][]string)
	fmt.Printf("handle pcl info%#v\n", pclInfo)
	lifdata := debug.GetLifData()

	MyMap := map[string][]LifRest{}
	for _, value := range lifdata.([]debug.D10LifEntry) {
		slice := MyMap[value.Server]
		element := LifRest{value.LifID, value.VlanId, value.PIF}
		slice = append(slice, element)
		MyMap[value.Server] = slice
	}

	sliceInfoRest := []D10LifInfoRest{}
	for server, lifinfo := range MyMap {
		sliceInfo := D10LifInfoRest{}
		sliceInfo.Server = server
		for _, v := range lifinfo {
			sliceInfo.LifSlice = append(sliceInfo.LifSlice, v)
		}
		sliceInfoRest = append(sliceInfoRest, sliceInfo)
	}
	fmt.Printf("lifdata Rest:%#v\n", sliceInfoRest)

	vlandata := debug.GetVlanFloodData()
	MyMapv := map[string][]D10VlanRest{}
	for _, value := range vlandata.([]debug.D10VlanFlood) {
		slice := MyMapv[value.Server]
		floodstr := fmt.Sprintf("%016b", value.LifFLoodMap)
		element := D10VlanRest{value.VlanID, floodstr}
		slice = append(slice, element)
		MyMapv[value.Server] = slice
	}
	slicevRest := []D10VlanInfoRest{}
	for server, vinfo := range MyMapv {
		sliceInfo := D10VlanInfoRest{}
		sliceInfo.Server = server
		for _, v := range vinfo {
			sliceInfo.D10Vlan = append(sliceInfo.D10Vlan, v)
		}
		slicevRest = append(slicevRest, sliceInfo)
	}

	fmt.Printf("lifdata Rest:%#v\n", slicevRest)

	w.Header().Add("Content Type", "text/html")
	templatesV := template.New("templateV")
	templatesV.New("Body").Parse(d10VDoc)
	templatesV.New("List").Parse(d10VDocList)
	templatesV.New("List1").Parse(d10VDocList1)

	templatesLif := template.New("templateLif")
	templatesLif.New("Body").Parse(lifDoc)
	templatesLif.New("List").Parse(lifDocList)
	templatesLif.New("List1").Parse(lifDocList1)

	templates := template.New("template")
	templates.New("Body").Parse(pclDoc)
	templates.New("List").Parse(pclDocList)
	templates.New("List1").Parse(pclDocList1)

	pclInfoRest := []PclInfoRest{}
	for server, ruleList := range pclInfo {
		pclrest := PclInfoRest{}
		pclrest.Server = server
		for _, rule := range ruleList {
			myrule := RuleListRest{rule}
			pclrest.RuleList = append(pclrest.RuleList, myrule)
		}
		pclInfoRest = append(pclInfoRest, pclrest)
	}

	fmt.Println("Lengh of lIF slice ", len(sliceInfoRest))
	fmt.Println("Lengh of PCL slice ", len(pclInfoRest))
	fmt.Println("Lengh of D10 Vlan list ", len(slicevRest))

	p := page{Title: "PCL VLAN Information", PclInfo: pclInfoRest}
	p1 := page{Title: "LIf Information", D10lifSlice: sliceInfoRest}
	p2 := page{Title: "Vlan Information", D10VlanSlice: slicevRest}
	if len(pclInfoRest) > 0 {
		templates.Lookup("Body").Execute(w, p)
	}
	if len(sliceInfoRest) > 0 {
		templatesLif.Lookup("Body").Execute(w, p1)
	}
	if len(sliceInfoRest) > 0 {
		templatesV.Lookup("Body").Execute(w, p2)
	}

}
func (r *RestObj) HandleMacInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetMacInfo()
	macinfo := val.(map[debug.MacVlan]map[debug.GPort]struct{})

	w.Header().Add("Content Type", "text/html")
	templates := template.New("template")
	templates.New("Body").Parse(macDoc)
	templates.New("List").Parse(macDocList)
	templates.New("List1").Parse(macDocList1)

	macInfoRest := []MacInfoRest{}
	for macvlan, serverportmap := range macinfo {
		smacinfo := MacInfoRest{}
		smacinfo.Mac = macvlan.Mac
		smacinfo.Vlan = macvlan.Vlan

		for portinfo, _ := range serverportmap {
			serverport := ServerPortInfo{Server: portinfo.Server,
				Port: portinfo.Portname}
			smacinfo.ServerPortSlice = append(smacinfo.ServerPortSlice, serverport)
		}
		macInfoRest = append(macInfoRest, smacinfo)

	}
	page := page{Title: "MAC Information", MacInfo: macInfoRest}
	templates.Lookup("Body").Execute(w, page)

}
func (r *RestObj) HandleErrInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetErrorInfo()
	myslice := val.([]debug.MyError)

	w.Header().Add("Content Type", "text/html")

	templates := template.New("template")
	templates.New("Body").Parse(errDoc)
	templates.New("List").Parse(errDocList)

	page := page{Title: "Error Information", ErrInfo: myslice}
	templates.Lookup("Body").Execute(w, page)

}

func (r *RestObj) HandleSystemInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetSystemInfo()
	mymap := val.(map[string]debug.SystemInfo)

	w.Header().Add("Content Type", "text/html")

	templates := template.New("template")
	templates.New("Body").Parse(sysDoc)
	templates.New("List").Parse(sysDocList)

	sysinfo := make([]debug.SystemInfo, 0)
	for _, element := range mymap {
		sysinfo = append(sysinfo, element)
	}
	fmt.Printf("slice rest %#v\n", sysinfo)

	page := page{Title: "System Information", SysInfo: sysinfo}
	templates.Lookup("Body").Execute(w, page)

}

func (r *RestObj) HandleResultsConfig(resp http.ResponseWriter, req *http.Request) {
	filepath := "config/results.html"
	http.ServeFile(resp, req, filepath)
}

func (r *RestObj) HandleDebugConfig(resp http.ResponseWriter, req *http.Request) {
	// Parse the form - this contains all the form elements
	// req.Form[<nameinhtml>][0] : This is map[string][]string
	req.ParseForm()

	path := req.Form["path"]
	servname := req.Form["servername"][0]
	cmd := req.Form["debugcmd"][0]

	fmt.Println("cmd is:", cmd)

	if cmd == "DebugClear" {
		resp.Write([]byte("<h1>Clearing all the DB</h1>"))
		debug.ClearDB()
		return
	}

	for _, val := range path {
		debug.Start(val, servname)
		// resp.Write([]byte("Path is:" + val))
		filepath := "config/results.html"
		http.ServeFile(resp, req, filepath)

	}
}
func (r *RestObj) HandleTestcases(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	setupinfo := req.Form["setup"][0]

	if len(setupinfo) == 0 {
		resp.Write([]byte("Setup Information Not given\n"))
		return
	}
	resp.Write([]byte(setupinfo + "\n"))

	for key, valueSlice := range req.Form {
		if strings.EqualFold(key, "testcase") {
			for index, value := range valueSlice {
				buf := fmt.Sprintf("Key:%s, Index:%d, Value:%s\n", key, index, value)
				resp.Write([]byte(buf))
			}
		}
	}
}
