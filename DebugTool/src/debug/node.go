package debug

import (
	"DebugTool/src/api"
	"DebugTool/src/unversioned"
	"encoding/json"
	"io/ioutil"
	"os"
)

// MgmtInterface contains the information of the Management interface on the
// server.
type MgmtInterface struct {
	Hostname  string   `json:"hostname"`
	IP        string   `json:"ip"`
	Gateway   string   `json:"gateway"`
	Interface string   `json:"interface"`
	MAC       string   `json:"mac"`
	DNS       []string `json:"dns"`
	Status    string   `json:"status"`
}

// OverlayInterface contains the information of the underline VF interface on the
// server used by overlay network.
type OverlayInterface struct {
	OverlayIP        string `json:"overlayIp"`
	OverlayInterface string `json:"overlayInterface"`
}

// Nodespec contains the specification of a node in the cluster.
type NodeSpec struct {
	Unschedulable bool `json:"unschedulable"`
	// Optional: Zone name
	Zone string `json:"zone,omitempty"`
}
type NodeState string
type NodeRole string

const (
	NodeRoleMaster NodeRole = "master"
	NodeRoleQuorum NodeRole = "etcd"
	NodeRoleWorker NodeRole = "worker"
)

// NodeStatus contains the status of a node in the cluster.
type NodeStatus struct {
	State NodeState `json:"state"`
	// Role is the role this node plays in the cluster: master, quorum, worker
	Role NodeRole `json:"role,omitempty"`
	// Optional: Schedulability of the node.
	Unschedulable bool `json:"unschedulable"`
	// Optional: Kubernetes Schedulability of the node.
	K8sUnschedulable bool `json:"k8sUnschedulable"`
	// Required: Last updated time
	LastUpdated string `json:"lastUpdated"`
	// Optional: message to indicate error
	Message string `json:"message,omitempty"`
	// Optional: Conditions of the node.
	NodeHealth api.ObjectHealth `json:"nodeHealth"`
	// Optional: Kubernetes node health status
	K8sHealth api.ObjectHealth `json:"k8sHealth"`
	// Optional: Indicates if node is recently rebooted.
	NeedsDiscovery bool `json:"needsDiscovery"`
}

// HostInfo contains the server host information.
type HostInfo struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:",inline"`
	MgmtInterface        `json:",inline"`
	OverlayInterface     `json:",inline"`
	IPMIAddress          string `json:"ipmiAddress"`
	IPMIMAC              string `json:"ipmiMAC"`
	IPMIGateway          string `json:"ipmiGateway"`
	OS                   string `json:"os"`
	KernelVersion        string `json:"kernelVersion"`
	ISOVersion           string `json:"isoVersion"`
	DockerVersion        string `json:"dockerVersion"`
	KubernetesVersion    string `json:"kubernetesVersion"`
	RecentBootTime       string `json:"recentBootTime"`
	Model                string `json:"serverModel"`
	CPU                  string `'json:"cpu"`
	SerialNumber         string `json:"serialNumber"`
	Runtime              string `json:"runtime"`
	RuntimeVersion       string `json:"runtimeVersion"`
	HostType             string `json: "hostType"`
}

type StorageMode string

const (
	L2StorageEnabled      StorageMode = "l2StorageEnabled"
	L3StorageEnabled      StorageMode = "l3StorageEnabled"
	L3StorageIpConfigured StorageMode = "l3StorageIpConfigured"
)

// StorageIPsInfo contains the storageIPs information.
type StorageIPsInfo struct {
	StorageIpPort1 string      `json:"storageIpPort1,omitempty"`
	StorageIpPort3 string      `json:"storageIpPort3,omitempty"`
	GatewayIP      string      `json:"gatewayIP,omitempty"`
	GatewayMac     string      `json:"gatewayMac,omitempty"`
	Subnet         string      `json:"subnet,omitempty"`
	SVLAN          uint        `json:"svlan,omitempty"`
	NetworkID      string      `json:"networkID,omitempty"`
	Mode           StorageMode `json:"mode,omitempty"`
}

// Node contains the status of a node in the cluster.
type Node struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:",inline"`
	Spec                 NodeSpec       `json:"spec,omitempty"`
	Status               NodeStatus     `json:"status,omitempty"`
	Host                 HostInfo       `json:"host,omitempty"`
	StorageIPs           StorageIPsInfo `json:"storageIPs,omitempty"`
}

// NodeList is a list of nodes.
type NodeList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:",inline"`

	Items     []Node `json:"items,omitempty"`
	GPUsCount int32  `json:"gpus_count,omitempty"`
}

func ProcessNodeInfo(path string, servername string) {

	f, _ := os.Open(path)

	b, _ := ioutil.ReadAll(f)

	var decodeData NodeList
	json.Unmarshal(b, &decodeData)
	//fmt.Printf("%#v\n", decodeData)
	writeToDb(decodeData)

}

func GetNodeInfo() interface{} {
	node := NodeList{}
	return readFromDb(node)
}
