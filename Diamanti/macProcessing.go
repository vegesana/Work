package main

import (
	"regexp"
	"strings"
)

// Index: 1256, Mac: e8:66:c4:7f:ff:f0, Vlan: 4094, Static: Yes, Port: 0
func MacInfoFun(line lineData) {
	macregexp := `([0-9a-f:]+[0-9a-f]+)`
	if strings.Contains(line.line, "Index") {
		Debug("String is::", line.line)
		r, _ := regexp.Compile(`Index: (\d+), Mac: ` + macregexp +
			`, Vlan: (\d+), Static: (\w+), Port: (\d+)`)
		if result := r.FindStringSubmatch(line.line); len(result) != 0 {
			macvlan := MacInfo{}
			//index := result[1]
			macvlan.mac = result[2]
			macvlan.vlan = result[3]
			//isStatic := result[4]
			macvlan.portNum = result[5]
			macvlan.fileName = line.filename
			Debug("Raju: PostMac:", macvlan)
			Gdata.GWriteCh <- macvlan
		}
	}

}
