package main

import (
	"regexp"
	"strings"
)

// Index: 1256, Mac: e8:66:c4:7f:ff:f0, Vlan: 4094, Static: Yes, Port: 0
func processMacInfo(line string, mydata Data) {
	if strings.Contains(line, "Index") {
		Debug("received data :", line)
		regexp.Compile("(?P<Index:>\\s\\d+)")
	}

}
