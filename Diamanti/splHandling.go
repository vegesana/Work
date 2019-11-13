package main

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func SplHandle() {
	cmd := exec.Command("/bin/sh", "myshell.sh")
	cmd.Run()
}
func splHandlingGo(tempDir string) {
	dir, _ := filepath.Glob(tempDir + "*.tmp")
	for _, filename := range dir {
		datab, _ := ioutil.ReadFile(filename)
		data := string(datab)
		Error("filename:", filename)
		name := filepath.Base(filename)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		name = strings.TrimSuffix(name, filepath.Ext(name))
		Error("rajuname:", name)

		r, _ := regexp.Compile(`Port: (\d+), Link: \w+, Pvid: \d+, MTU: \d+, Lane0_Tx_Polarity:\s+` +
			`\w+, Lane0_Rx_Polarity: \w+\s+Protected_Mode: \w+, dstPort: (\d+), dstHwDev: \d+`)
		if result := r.FindAllStringSubmatch(string(data), -1); len(result) != 0 {
			for _, element := range result {
				// we need this 2 times to trip two extensions
				sport := GPort{name, element[1]}
				dport := GPort{name, element[2]}
				pininfo := PinInfo{sport, dport}
				Error("RajuPost", pininfo)
				Gdata.GWriteCh <- pininfo
			}
		}
	}
}
