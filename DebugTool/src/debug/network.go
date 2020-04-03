package debug

import (
	"DebugTool/src/api"
	"DebugTool/src/unversioned"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type NetworkSpec struct {
	Type string `json:"type,omitempty"`

	Subnet string `json:"subnet"`

	StartAddr string `json:"startAddr,omitempty"`

	EndAddr string `json:"endAddr,omitempty"`

	GatewayIP string `json:"gatewayIP,omitempty"`

	VLAN uint `json:"vlan,omitempty"`

	Routes []string `json:"routes,omitempty"`

	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

type Network struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:",inline"`

	// Spec specifies configuration information of the network.
	Spec NetworkSpec `json:"spec,omitempty"`

	// StorageNetwork specifies configuration information of the network.
	StorageNetworkSpec StorageNWSpec `json:"storageNetworkSpec,omitempty"`
	Status             NetworkStatus `json:"status,omitempty"`
}

type NetworkStatus struct {
	Message string `json:"message,omitempty"`

	Subnet string `json:"subnet"`

	// Start of the allocatable addresses in the Network.
	StartAddr string `json:"startAddr"`

	// IP address of the gateway for the Network.
	GatewayIP string `json:"gatewayIP,omitempty"`

	// VLAN/VNI id of network.
	VLAN uint `json:"vlan,omitempty"`

	// Number of Used IP addresses from the pool
	UsedAddrs int `json:"usedAddrs"`

	// Total IP addresses in the pool
	NumAddrs int `json:"numAddrs"`

	// Number of attached containers
	NumContainers int `json:"numContainers,omitempty"`

	// Administrative Endpoint State.
	MarkedForDeletion bool `json:"markedForDeletion,omitempty"`
}
type StorageNWSpec struct {
	// Required: gateway mac for storage-network.
	GatewayMac string `json:"gatewayMac,omitempty"`
}

type OverlayNetworkConfig struct {
	Name             string
	Network          string
	HostSubnetLength uint32
	ServiceNetwork   string
	PluginName       string
}
type NetworkList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:",inline"`
	Items                []Network `json:"items,omitempty"`
}

func ProcessNetworkInfo(path string, servername string) {

	f, _ := os.Open(path)

	b, _ := ioutil.ReadAll(f)

	var decodeData NetworkList
	json.Unmarshal(b, &decodeData)
	fmt.Printf("Total Item : %d\n", decodeData.TotalItems)
	// fmt.Printf("%#v\n", decodeData)
	writeToDb(decodeData)

}

func GetNetworkInfo() interface{} {
	networkinfo := NetworkList{}
	return readFromDb(networkinfo)
}
