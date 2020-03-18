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
	ErrInfo []debug.MyError
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
func (r *RestObj) HandleMacInfo(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Handle Mac Info"))

}
func (r *RestObj) HandleErrInfo(w http.ResponseWriter, req *http.Request) {
	val := debug.GetErrorInfo()
	myslice := val.([]debug.MyError)

	w.Header().Add("Content Type", "text/html")

	templates := template.New("template")
	templates.New("Body").Parse(errDoc)
	templates.New("List").Parse(errDocList)

	fmt.Printf("slice rest %#v\n", myslice)

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

	for _, val := range path {
		debug.Start(val)
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

const macDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.ServerName}}</td>
        <td>{{.BoardInfo}}</td>
        <td>{{.ProductId}}</td>
    </tr>
    {{end}}
</ul>
`

const macDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
        <table style="width:20%" border="1">
        <tr>
            <th> ServerName </th>
            <th> BoardInfo</th>
            <th> ProductId</th>
        </tr>
        {{template "List" .SysInfo}}
        </table>
    </body>
</html>
`
const errDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.ServerName}}</td>
        <td>{{.MyErr}}</td>
    </tr>
    {{end}}
</ul>
`

const errDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
        <table style="width:40%" border="1">
        <tr>
            <th> ServerName </th>
            <th> Error </th>
        </tr>
        {{template "List" .ErrInfo}}
        </table>
    </body>
</html>
`
const sysDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.ServerName}}</td>
        <td>{{.BoardInfo}}</td>
        <td>{{.ProductId}}</td>
    </tr>
    {{end}}
</ul>
`

const sysDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
        <table style="width:20%" border="1">
        <tr>
            <th> ServerName </th>
            <th> BoardInfo</th>
            <th> ProductId</th>
        </tr>
        {{template "List" .SysInfo}}
        </table>
    </body>
</html>
`
