package debug

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type D10LifEntry struct {
	Server string
	LifID  uint64
	VlanId uint64
	PIF    uint64
}
type D10VlanFlood struct {
	Server      string
	VlanID      uint64
	LifFLoodMap uint64
}

func processL2Util(path, servername string) error {
	f, _ := os.Open(path)
	b, _ := ioutil.ReadAll(f)

	//	var d10LifDb map[string]
	expr := `LIF entry\s+(\d+)\s+VNID\s+(\d+)\s+VLAN id\s+` +
		`(\w+)\s+VxLAN LIF flag\s+(\d+)\s+Pinned NIF\s+(\d+)\s+` +
		`PIF index\s+(\d+)\s+VLAN rewrite flag\s+(\d+)\s+LIF valid flag\s+1`
	r, _ := regexp.Compile(expr)
	sliceSlice := r.FindAllStringSubmatch(string(b), -1)
	for _, slice := range sliceSlice {
		lifid := slice[1]
		ilifid, _ := strconv.ParseUint(lifid, 16, 32)
		vlanid := slice[3]
		vvid, _ := strconv.ParseUint(vlanid, 16, 32)
		pifvalue := slice[6]
		pInt, _ := strconv.ParseUint(pifvalue, 16, 32)
		fmt.Printf("lifid :%d vlan :%d, pif:%d \n", ilifid, vvid, pInt)
		writeToDb(D10LifEntry{servername, ilifid, vvid, pInt})

	}

	expr = `VLAN id\s+(\w+)\s+VLAN valid flag\s+1\s+` +
		`LIF flood map bitmap\s+(\d+)\s+`
	r, _ = regexp.Compile(expr)
	sliceSlice = r.FindAllStringSubmatch(string(b), -1)

	vlanFlood := make(map[uint64]uint64)
	for _, slice := range sliceSlice {
		vlanid := slice[1]
		vvid, _ := strconv.ParseUint(vlanid, 16, 32)
		flood := slice[2]
		fvalue, _ := strconv.ParseUint(flood, 16, 32)
		vlanFlood[vvid] = fvalue
	}

	expr = `FTBL\s+(\w+)\s+LIF flood table\s+(\w+ \w+ \w+ \w+)`
	r, _ = regexp.Compile(expr)
	sliceSlice = r.FindAllStringSubmatch(string(b), -1)

	floodStr := make(map[uint64]uint64)
	for _, slice := range sliceSlice {
		fldtbl := slice[1]
		fInt, _ := strconv.ParseUint(fldtbl, 16, 32)
		flood := slice[2]
		newstr := strings.Replace(flood, " ", "", -1)
		fvalue, _ := strconv.ParseUint(newstr, 16, 64)
		if fvalue != 0 {
			floodStr[fInt] = fvalue
		}
	}

	for vid, floodid := range vlanFlood {
		writeToDb(D10VlanFlood{servername, vid, floodStr[floodid]})
	}
	return nil
}
func GetLifData() interface{} {
	info := D10LifEntry{}
	return readFromDb(info)
}
func GetVlanFloodData() interface{} {
	info := D10VlanFlood{}
	return readFromDb(info)
}
