package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Key struct {
	lport, session, key string
}

// valid session
type VSession struct {
	lport   string
	session string
}

var btpUtilMap map[Key]string
var AllStringKeyMap map[string]struct{}

func DumpBtpErrors(servername string) {

	slice := getAllValidSessions()
	fmt.Println("Valid Sessions", slice)

	for _, value := range slice {
		for str, _ := range AllStringKeyMap {
			k := Key{value.lport, value.session, str}
			v := btpUtilMap[k]
			if strings.Contains(k.key, "Number of queued NVMEoE command frames") {
				value, _ := strconv.Atoi(v)
				if value != 0 {
					reason := "NVMEoE Commands are queued and not dequed by BTP"
					err := fmt.Sprintf("lport:%s,sess:%s,val:%d: %s\n",
						k.lport, k.session, value, reason)
					SendError(err)
				}
			}
			if strings.Contains(k.key, "Number of queued NVMEoE control frames") {
				value, _ := strconv.Atoi(v)
				if value != 0 {
					reason := "NVMEoE COntrol are queued and not dequed by BTP"
					err := fmt.Sprintf("lport:%s,sess:%s,val:%d: %s\n",
						k.lport, k.session, value, reason)
					SendError(err)

				}
			}
			if strings.Contains(k.key, "Number of queued NVMEoE data frames") {
				value, _ := strconv.Atoi(v)
				if value != 0 {
					reason := "NVMEoE Data are queued and not dequed by BTP"
					err := fmt.Sprintf("lport:%s,sess:%s,val:%d: %s\n",
						k.lport, k.session, value, reason)
					SendError(err)
				}
			}
			if strings.Contains(k.key, "BTP ping FSM state") {
				value, _ := strconv.Atoi(v)
				if value != 1 {
					/*
						1: Initally and after we get response
							we will move to this state
						2: Ping sent but waiting for response
					*/
					reason := "Session Ping not in Good state"
					err := fmt.Sprintf("%s:<lport:%s,sess:%s>, val:%d: %s\n",
						servername, k.lport, k.session, value, reason)
					SendError(err)

				}
			}

			if strings.Contains(k.key, "RS login state") {
				value, _ := strconv.Atoi(v)
				if value != 1 {
					reason := "Session Rslogin not in Good state"
					err := fmt.Sprintf("%s:<lport:%s,sess:%s>, val:%d: %s\n",
						servername, k.lport, k.session, value, reason)
					SendError(err)

				}
			}

		}
	}
	// Check The BTP session
	// ping_state : BTP_PING_ST_REQ_SENT (sent and waiting for response)
}

func processBtpResult() {

}
func getAllValidSessions() []VSession {
	myslice := make([]VSession, 0)

	for k, v := range btpUtilMap {
		if strings.Contains(k.key, "Session State") {
			value, _ := strconv.Atoi(v)
			if value == 1 {
				myslice = append(myslice, VSession{k.lport, k.session})
			}
		}
	}
	return myslice
}
func processBtpUtil(path string, servername string) error {
	// <lportid><sessionid><key>:value
	btpUtilMap = map[Key]string{}
	AllStringKeyMap = map[string]struct{}{}

	var lstart, sstart bool
	var lport, session string
	f, _ := os.Open(path)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`Lport (\d+)+`)
		lslice := re.FindStringSubmatch(line)
		if len(lslice) > 1 {
			lstart = true
			lport = lslice[1]
			continue
		}
		if lstart {
			re := regexp.MustCompile(`Session (\d+)+`)
			lslice := re.FindStringSubmatch(line)
			if len(lslice) > 1 {
				sstart = true
				session = lslice[1]
				continue
			}
			if lstart && sstart {
				if strings.Contains(line, ":") {
					pair := strings.Split(line, ":")
					mystr := strings.TrimSpace(pair[0])
					if _, ok := AllStringKeyMap[mystr]; !ok {
						AllStringKeyMap[mystr] = struct{}{}
					}
					key := Key{lport, session, mystr}
					btpUtilMap[key] = strings.TrimSpace(pair[1])
				}
			}

		}

	}
	DumpBtpErrors(servername)
	return nil
}
