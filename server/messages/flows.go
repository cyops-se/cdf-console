package messages

// timestamp: 2025-01-31 20:49:50 +0000 UTC
// hostname: RUTX08
// client: 10.49.88.8:40822
// content: [75281.561206] WGHUB-GREY IN=vxlan3 OUT=vxlan1 MAC source = dc:a6:32:35:bf:50 MAC dest = 01:00:5e:00:00:fb proto = 0x0800 IP SRC=172.16.91.32 IP DST=224.0.0.251, IP tos=0x00, IP proto=17 SPT=5353 DPT=5353

type FlowEntry struct {
	Hostname  string `json:"hostname"` // Reporting hostname
	IpAddress string `json:"ip"`       // Reporting IP address
	Prefix    string `json:"prefix"`
	IfaceIn   string `json:"ifacein"`
	IfaceOut  string `json:"ifaceout"`
	HwProto   string `json:"hwproto"`
	SrcMac    string `json:"srcmac"`
	DstMac    string `json:"dstmac"`
	IpProto   string `json:"ipproto"`
	SrcIp     string `json:"srcip"`
	SrcPort   string `json:"srcport"`
	DstIp     string `json:"dstip"`
	DstPort   string `json:"dstport"`
	Count     uint   `json:"count"`
}

type EntryListResponse struct {
	RequestResponse
	EntriesList []*FlowEntry `json:"entries"`
}

type EtherType struct {
	Value       string `json:"value"` // string in C hex form; 0x0800
	Description string `json:"description"`
}

type EtherTypeListResponse struct {
	RequestResponse
	EtherTypeList []EtherType `json:"ethertypes"`
}
