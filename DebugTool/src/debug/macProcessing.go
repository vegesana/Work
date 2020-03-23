package debug

import (
	"regexp"
	"strings"
)

// Index: 1256, Mac: e8:66:c4:7f:ff:f0, Vlan: 4094, Static: Yes, Port: 0
func MacInfoFun(line lineData) {
	macregexp := `([0-9a-f:]+[0-9a-f]+)`
	if strings.Contains(line.Line, "Index") {
		Debug("String is::", line.Line)
		r, _ := regexp.Compile(`Index: (\d+), Mac: ` + macregexp +
			`, Vlan: (\d+), Static: (\w+), Port: (\d+)`)
		if result := r.FindStringSubmatch(line.Line); len(result) != 0 {
			macvlan := MacInfo{}
			//index := result[1]
			macvlan.Mac = result[2]
			macvlan.Vlan = result[3]
			//isStatic := result[4]
			macvlan.PortNum = result[5]
			macvlan.Server = line.Server
			writeToDb(macvlan)
		}
	}
}
