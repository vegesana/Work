package rest

import (
	"fmt"
	"net/http"
)

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
	go http.ListenAndServe("Localhost:8080", nil)
}

func (r *RestObj) HandleMainConfig(resp http.ResponseWriter, req *http.Request) {
	filepath := "config/main.html"
	http.ServeFile(resp, req, filepath)
}
