package main

var temp string

func processPinInfo(line lineData) {
	Debug("PIN String is::", line.line)
	if line.line == "\n" {
		Debug("Raju: Pin Processing", temp)
		temp = ""
	} else {
		temp = temp + line.line
	}
}
