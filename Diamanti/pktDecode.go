package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Pkt struct {
	damac    string
	samac    string
	priority string
	cfi      string
	vlan     string
	pkttype  string
	ipkttype string // Inner pkt type
	sip      string
	dip      string
}

func processFile() {
	pktslice := make([]Pkt, 0)
	fd, _ := os.Open("InputFiles/pktinput.txt")
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println("LIne is ", ln)
		if len(ln) > 1500 {
			pktslice = append(pktslice, decode(strings.TrimSpace(ln)))
		}

	}
	printPkt(pktslice)

}
func main() {

	processFile()

	watcher, _ := fsnotify.NewWatcher()
	watcher.Add("InputFiles/")

	for {
		select {
		case event := <-watcher.Events:
			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				fmt.Println("Create new file")
			case event.Op&fsnotify.Write == fsnotify.Write:
				fmt.Println("Wrote to this file")
				processFile()
			default:
				fmt.Println("Event happend", event.Op)
			}
		case errors := <-watcher.Errors:
			fmt.Println("Errors ", errors)
		}
	}
}

func printPkt(pktSlice []Pkt) {
	for i, pkt := range pktSlice {
		fmt.Println("Pkt:", i)
		fmt.Println("============")
		fmt.Printf("MACDA:%s, MACSA:%s, type:%s,",
			pkt.damac, pkt.samac, pkt.pkttype)
		if pkt.pkttype == "8100" {
			vlanid, _ := strconv.ParseInt(pkt.vlan, 16, 64)
			pricfi, _ := strconv.ParseInt(pkt.priority, 16, 32)
			pri := (pricfi >> 1) & 0xF
			cfi := (pricfi) & 0x1
			fmt.Printf("pri:%d, cfi:%d vlan:%d", pri, cfi, vlanid)
		}

		if pkt.ipkttype == "88B7" {
			fmt.Printf(",BTP Packet")
		} else {
			fmt.Printf("NON BTP Packet", pkt.ipkttype)
		}
		fmt.Println("\n---------------")
	}
}
func decode(mystr string) Pkt {
	var pkt Pkt
	slice := strings.Split(mystr, "")
	fmt.Println("Lenght of slice ", len(slice))
	start := 0
	end := start + 12
	pkt.damac = MacAddress(slice[start:end])
	start = end
	end = start + 12
	pkt.samac = MacAddress(slice[start:end])
	start = end
	end = start + 4
	pkt.pkttype = PktType(slice[start:end])
	start = end
	if pkt.pkttype == "8100" {
		end = start + 4
		pkt.priority = slice[start]

		pkt.vlan = strings.Join(slice[start+1:end], "")
		start = end
	}
	// Pkt type - 0800 or 0x88B7 (BTP)
	end = start + 4
	pkt.ipkttype = PktType(slice[start:end])

	return pkt
}

func PktType(pkttype []string) string {
	return strings.Join(pkttype, "")
}
func MacAddress(mac []string) string {
	var substr string
	var newslice []string

	fmt.Println("Input Mac address", strings.Join(mac, ""))
	for i := 1; i <= 12; i++ {
		substr = substr + mac[i-1]
		if i%2 == 0 {
			newslice = append(newslice, substr)
			substr = ""
		}
	}
	return strings.Join(newslice, ":")
}
