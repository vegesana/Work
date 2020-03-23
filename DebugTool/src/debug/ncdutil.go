package debug

import (
	"bufio"
	"os"
	"strings"
)

func processNcdUtil(path string, servername string) error {

	f, _ := os.Open(path)
	go func(f *os.File, servername string, mych chan interface{}) {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			origTxt := scanner.Text()
			ln := strings.TrimSpace(origTxt)

			if strings.Index(ln, "Network Interface Info") == 0 {
				Debug("Picked Network Interface Channel")
				mych = NetworkInfoCh
			}

			if strings.Index(ln, "MAC Info") == 0 {
				Debug("MAC Info channel ")
				mych = MacInfoCh
			}

			if strings.Index(ln, "Interface Info") == 0 {
				Debug("Interface Info channel ")
				mych = IntfInfoCh
			}

			if strings.Index(ln, "PCL Info") == 0 {
				Debug("PCL Info channel ")
				mych = PclInfoCh
			}

			lninfo := lineData{servername, origTxt, servername}
			mych <- lninfo

			// All the end to determine the end of Channel
			if strings.Contains(ln, "maxPortDescrLimit") {
				// Place HOlder till we find new channel
				mych = JunkCh
			}

		}
		if err := scanner.Err(); err != nil {
			Debug("Err:", err.Error())
		}
	}(f, servername, JunkCh)
	return nil
}

func goRoutine(str string, ch chan interface{}, fun func(lineData)) {

	for {
		select {
		case val := <-ch:
			switch val.(type) {
			case lineData:
				myval := val.(lineData)
				// Check wheather we have any function pointer
				// for the given
				if fun != nil {
					fun(myval)
				} else {
					Debug("NO processing defined for :")
				}

			default:
			}

		}
	}
}

func (line lineData) getText() string {
	return line.Line
}
func (line lineData) getFileName() string {
	return line.Server
}
