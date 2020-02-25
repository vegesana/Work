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
	odmac      string
	osmac      string
	priority   string
	cfi        string
	vlan       string
	ivlan      string
	vni        string
	pkttype    string
	ipkttype   string // Inner pkt type
	iipkttype  string // Inner pkt type
	iiipkttype string // Inner pkt type
	osip       string
	odip       string
	sport      string
	dport      string
	ipcksum    string
	idmac      string
	ismac      string
	isip       string
	idip       string
}

func processFile() {
	pktslice := make([]Pkt, 0)
	fd, _ := os.Open("InputFiles/pktinput.txt")
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println("LIne is ", ln)
		if len(ln) > 200 {
			pktslice = append(pktslice, decode(strings.TrimSpace(ln)))
		} else {
			fmt.Println("Lenght of packet is less than 200")
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
			pkt.odmac, pkt.osmac, pkt.pkttype)
		if pkt.pkttype == "8100" {
			vlanid, _ := strconv.ParseInt(pkt.vlan, 16, 64)
			pricfi, _ := strconv.ParseInt(pkt.priority, 16, 32)
			pri := (pricfi >> 1) & 0xF
			cfi := (pricfi) & 0x1
			fmt.Printf("pri:%d, cfi:%d vlan:%d\n", pri, cfi, vlanid)
		}

		if pkt.ipkttype == "88B7" {
			fmt.Println("BTP Packet")
		} else if pkt.ipkttype == "0800" {
			fmt.Printf("IP Packet")
			fmt.Printf(",IPCKSUM :%s", pkt.ipcksum)
			fmt.Printf(",Outer DIP:%s", pkt.odip)
			fmt.Printf(",Outer SIP::%s", pkt.osip)
			if pkt.dport == "12B5" {
				fmt.Printf(",VXLAN UDP Packet")
				fmt.Printf(",VNI 0x%s", pkt.vni)
				fmt.Printf(",IDMAC:%s", pkt.idmac)
				fmt.Printf(",ISMAC:%s", pkt.ismac)
				fmt.Printf(",INner PktType:%s", pkt.iipkttype)

				vlanid, _ := strconv.ParseInt(pkt.ivlan, 16, 64)
				fmt.Printf(",INner Vlan :%d", vlanid)
				fmt.Printf(",Acutual Pkt type:%s", pkt.iiipkttype)
			}
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
	end := 12
	pkt.odmac = MacAddress(slice[start:end])
	start = 12
	end = 24
	pkt.osmac = MacAddress(slice[start:end])
	start = 24
	end = 28
	pkt.pkttype = PktType(slice[start:end])
	start = 28
	if pkt.pkttype == "8100" {
		end = start + 4
		pkt.priority = slice[start]

		pkt.vlan = strings.Join(slice[start+1:end], "")
		start = end
	}
	// Pkt type - 0800 or 0x88B7 (BTP)
	end = start + 4
	pkt.ipkttype = PktType(slice[start:end])

	if pkt.ipkttype == "0800" {

		start = 56
		end = start + 4
		pkt.ipcksum = Port(slice[start:end])

		start = 60
		end = start + 8
		pkt.osip = IPAddress(slice[start:end])
		start = end
		end = start + 8
		pkt.odip = IPAddress(slice[start:end])
		start = end
		end = start + 4
		pkt.sport = Port(slice[start:end])
		fmt.Println("SPort is ", pkt.sport)

		start = end
		end = start + 4
		pkt.dport = Port(slice[start:end])
		fmt.Println("DPort is ", pkt.dport)

		if pkt.dport == "12B5" {
			start = end + 16
			end = start + 6
			pkt.vni = Port(slice[start:end])
			fmt.Println("pkt vni ", pkt.vni)
			start = end + 2
			end = start + 12
			pkt.idmac = MacAddress(slice[start:end])
			start = end
			end = start + 12
			pkt.ismac = MacAddress(slice[start:end])
			start = end
			end = start + 4
			pkt.iipkttype = PktType(slice[start:end])
			start = end

			if pkt.iipkttype == "8100" {
				end = start + 4
				pkt.priority = slice[start]

				pkt.ivlan = strings.Join(slice[start+1:end], "")
				start = end
			}
			end = start + 4
			pkt.iiipkttype = PktType(slice[start:end])

		}
	}

	return pkt
}
func Port(port []string) string {
	return strings.Join(port, "")
}
func PktType(pkttype []string) string {
	return strings.Join(pkttype, "")
}

func convertHexToInterSlice(slice []string) []string {
	var newslice []string
	for _, val := range slice {
		newval, _ := strconv.ParseUint(val, 16, 32)
		// fmt.Println("IP:", newval)
		newslice = append(newslice, fmt.Sprint(newval))

	}
	return newslice

}

func IPAddress(ip []string) string {
	var substr string
	var newslice []string
	ipstr := strings.Join(ip, "")
	fmt.Println("IP address is ", ipstr)
	for i := 1; i <= 8; i++ {
		substr = substr + ip[i-1]
		if i%2 == 0 {
			newslice = append(newslice, substr)
			substr = ""
		}
	}
	nnslice := convertHexToInterSlice(newslice)
	myip := strings.Join(nnslice, ".")
	fmt.Println("My ip is ", myip)
	return myip
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
