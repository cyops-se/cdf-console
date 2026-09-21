package types

import (
	"time"
)

type Network struct {
	Model
	Name     string `json:"name"`
	Network  string `json:"network"`
	Comment  string `json:"comment"`
	SystemID uint   `json:"systemid"`
}

type System struct {
	Model
	Name string   `json:"name"`
	Hubs []Device `json:"hubs"`
}

type WireguardPeer struct {
	Model
	Name                string `json:"name"`
	DeviceID            uint   `json:"deviceid"`
	Device              Device `json:"device"`     // Use device.wg.publickey for peer
	AllowedIps          string `json:"allowedips"` // comma separated list of ips
	WireguardInstanceID uint   `json:"wgid"`
}

type WireguardInstance struct {
	Model
	Name       string          `json:"name"`
	Ip         string          `json:"ip"`
	PublicKey  string          `json:"pubkey"`
	PrivateKey string          `json:"privkey"` // should be hashed for comparison, not in plain text
	DeviceID   uint            `json:"deviceid"`
	Peers      []WireguardPeer `json:"peers"`
}

type Address struct {
	Model
	DeviceID uint   `json:"deviceid"`
	Address  string `json:"address"`
}

// DeviceInterface represents a network interface with its configuration and statistics
// This is now a proper GORM model with its own table
// type DeviceInterface struct {
// 	Model
// 	DeviceID   uint           `json:"deviceid" gorm:"index"`
// 	Name       string         `json:"name" gorm:"index"` // 'wan', 'lan', 'apn'
// 	IpAddress  string         `json:"ipaddress"`
// 	Netmask    string         `json:"netmask"`
// 	MacAddress string         `json:"macaddress"`
// 	Present    bool           `json:"present"`
// 	Stats      InterfaceStats `json:"stats" gorm:"type:text"`
// }

// type WireguardInterface struct {
// 	Model
// 	DeviceID   uint   `json:"deviceid"`
// 	Name       string `json:"name"`       // Interface name (e.g., 'wg0', 'wg1')
// 	IpAddress  string `json:"ipaddress"`  // Wireguard tunnel IP
// 	Netmask    string `json:"netmask"`    // Wireguard tunnel netmask
// 	ListenPort uint   `json:"listenport"` // Wireguard listen port
// 	PublicKey  string `json:"publickey"`  // Wireguard public key
// 	PrivateKey string `json:"privatekey"` // Wireguard private key (encrypted)
// 	Endpoint   string `json:"endpoint"`   // Remote endpoint (for clients)
// 	AllowedIPs string `json:"allowedips"` // Allowed IPs for this interface
// 	Status     string `json:"status"`     // 'unknown', 'not configured', 'not active', 'active'
// }

// InterfaceStats holds traffic statistics for an interface
type InterfaceStats struct {
	Model
	DeviceID  uint   `json:"devid"`
	RxBytes   uint64 `json:"rxbytes"`
	TxBytes   uint64 `json:"txbytes"`
	RxPackets uint64 `json:"rxpackets"`
	TxPackets uint64 `json:"txpackets"`
	Operstate string `json:"operstate"`
}

type Device struct {
	Model
	SystemID        uint   `json:"systemid"`
	Role            string `json:"role"` // 'hub' | 'endpoint'
	SerialNo        string `json:"serialno"`
	DeviceModel     string `json:"devicemodel"`
	Firmware        string `json:"firmware"`
	Name            string `json:"name"`
	Hostname        string `json:"hostname"`
	DefaultPassword string `json:"defaultpwd"`
	Password        string `json:"password"`
	Postfix         string `json:"postfix"`

	// Legacy columns for backward compatibility and migration
	Wan          string         `json:"wan"`
	WanMask      string         `json:"wanmask"`
	WanPresent   bool           `json:"wanpresent"`
	WanStats     InterfaceStats `json:"wanstats" gorm:"foreignkey:DeviceID"`
	Apn4G        string         `json:"apn4g"`
	Apn4Gmask    string         `json:"apn4gmask"`
	Apn4GPresent bool           `json:"apn4gpresent"`
	Apn4GStats   InterfaceStats `json:"apn4gstats" gorm:"foreignkey:DeviceID"`
	Lan          string         `json:"lan"`
	Lanmask      string         `json:"lanmask"`
	LanPresent   bool           `json:"lanpresent"`
	LanStats     InterfaceStats `json:"lanstats" gorm:"foreignkey:DeviceID"`

	// WireguardInterface WireguardInterface `json:"wireguardinterface" gorm:"foreignkey:DeviceID"`
	PreferredIp string `json:"preferredip"`

	Version      string    `json:"version"`
	State        string    `json:"state"`        // 'pending' | 'adopted' | 'provisioned'
	Status       string    `json:"status"`       // 'unknown', 'not reachable', 'ok'
	SerialStatus string    `json:"serialstatus"` // 'unknown', 'mismatch', 'ok'
	SshStatus    string    `json:"sshstatus"`    // 'unknown', 'no keys', 'ok'
	WgStatus     string    `json:"wgstatus"`     // 'unknown', 'not configured', 'not active', 'active'
	LastCheck    time.Time `json:"lastcheck"`
	HubID        *uint     `json:"hubid"`
	Endpoints    []Device  `json:"endpoints" gorm:"foreignkey:HubID"`

	// Legacy wireguard field (TODO: migrate to WireguardInterfaces)
	PublicKey string `json:"pubkey"`

	AllowedAddresses []Address `json:"allowedaddresses"`
	SyslogIp         string    `json:"syslogip"`
	SyslogPort       string    `json:"syslogport"`
}
