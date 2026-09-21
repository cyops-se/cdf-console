package types

import "time"

type DataFlow struct {
	Model
	Protocol       string    `json:"protocol"`
	SrcMac         string    `json:"srcmac"`
	SrcIp          string    `json:"srcip"`
	SrcPort        int       `json:"srcport"`
	DstMac         string    `json:"dstmac"`
	DstIp          string    `json:"dstip"`
	DstPort        int       `json:"dstport"`
	Count          int       `json:"count"`
	LastSeen       time.Time `json:"lastseen"`
	DataFlowListID uint      `json:"dataflowlistid"`
}

type DataFlowList struct {
	Model
	Name      string     `json:"name"`
	ListType  string     `json:"listtype"`
	Allow     bool       `json:"allow"`
	Log       bool       `json:"log"`
	Notify    bool       `json:"notify"`
	DataFlows []DataFlow `json:"dataflows"`
}
