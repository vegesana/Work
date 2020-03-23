package rest

import (
	"DebugTool/src/debug"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type page struct {
	Title   string
	SysInfo []debug.SystemInfo
	MacInfo []MacInfoRest
	PinInfo []PinInfoRest
	PclInfo []PclInfoRest
	ErrInfo []debug.MyError
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

	w.Header().Add("Content Type", "text/html")
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

	page := page{Title: "PCL VLAN Information", PclInfo: pclInfoRest}
	templates.Lookup("Body").Execute(w, page)

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

	for _, val := range path {
		debug.Start(val, servname)
		// resp.Write([]byte("Path is:" + val))
		filepath := "config/results.html"
		http.ServeFile(resp, req, filepath)

	}
}
func (r *RestObj) HandleTestcases(resp http.ResponseWriter, req *http.Request) {
	// req.Form[<nameinhtml>][]string : This is map[string][]string
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
