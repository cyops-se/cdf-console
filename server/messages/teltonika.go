package messages

import "server/types"

// REQUESTS

// POST http://192.168.1.1/api/login
// {"username":"admin","password":"admin01"}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// {"data":{"password":"Hemligt01","password_confirm":"Hemligt01"}}
type ChangeDefaultPasswordRequestData struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ChangeDefaultPasswordRequest struct {
	Data ChangeDefaultPasswordRequestData `json:"data"`
}

type DeviceInfoRequest struct {
	Token      string `json:"token"`
	Hostname   string `json:"hostname"`
	DeviceName string `json:"devicename"`
	LanIp      string `json:"lanip"`
	LanMask    string `json:"lanmask"`
	WanIp      string `json:"wanip"`
	WanMask    string `json:"wanmask"`
}

type RunScriptRequest struct {
	Script    string `json:"script"`
	Arguments string `json:"args"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RunScript2Request struct {
	Script    string     `json:"script"`
	Arguments []KeyValue `json:"args"`
}

// RESPONSES

type LoginResponseData struct {
	Username string   `json:"username"`
	Token    string   `json:"token"`
	Expires  int      `json:"expires"`
	Errors   []string `json:"errors"`
}

type LoginRequestResponse struct {
	Success bool              `json:"success"`
	Data    LoginResponseData `json:"data"`
}

type RequestResponse struct {
	Success bool          `json:"success"`
	Errors  []interface{} `json:"errors"`
}

type DeviceInfoResponse struct {
	Success bool          `json:"success"`
	Data    types.Device  `json:"data"`
	Errors  []interface{} `json:"errors"`
}

/*
{
    "success": true,
    "data": {
        "mnfinfo": {
            "macEth": "001E425C5150",
            "name": "RUTX0800XXXX",
            "blver": "2.3",
            "hwver": "0010",
            "batch": "0024",
            "serial": "1125626527",
            "mac": "001E425C514F"
        },
        "static": {
            "fw_version": "RUTX_R_00.07.10.2",
            "kernel": "5.10.224",
            "system": "ARMv7 Processor rev 5 (v7l)",
            "device_name": "RUTX08-dev",
            "cpu_count": 4,
            "hostname": "RUTX08-host",
            "release": {
                "distribution": "OpenWrt",
                "revision": "r16279-5cc0535800",
                "version": "21.02.0",
                "target": "ipq40xx/generic",
                "description": "OpenWrt 21.02.0 r16279-5cc0535800"
            },
            "fw_build_date": "2024-10-30 10:59:57",
            "model": "RUTX08",
            "board_name": "teltonika,rutx"
        },
        "features": {
            "ipv6": true
        },
        "board": {
            "serial": [
                {
                    "external_devices": []
                }
            ],
            "network": {
                "wan": {
                    "device": "eth1",
                    "proto": "dhcp"
                },
                "lan": {
                    "device": "eth0",
                    "proto": "static",
                    "default_ip": "192.168.1.1"
                }
            },
            "model": {
                "id": "teltonika,rutx",
                "platform": "RUTX",
                "name": "RUTX08"
            },
            "switch": {
                "switch0": {
                    "enable": true,
                    "roles": [
                        {
                            "ports": "2 3 4 0",
                            "role": "lan",
                            "device": "eth0"
                        },
                        {
                            "ports": "5 0",
                            "role": "wan",
                            "device": "eth1"
                        }
                    ],
                    "ports": [
                        {
                            "device": "eth0",
                            "num": 0,
                            "want_untag": true,
                            "need_tag": false
                        },
                        {
                            "num": 2,
                            "index": 1,
                            "role": "lan"
                        },
                        {
                            "num": 3,
                            "index": 2,
                            "role": "lan"
                        },
                        {
                            "num": 4,
                            "index": 3,
                            "role": "lan"
                        },
                        {
                            "device": "eth1",
                            "num": 0,
                            "want_untag": true,
                            "need_tag": false
                        },
                        {
                            "role": "wan",
                            "num": 5
                        }
                    ],
                    "reset": true
                }
            },
            "network_options": {
                "readonly_vlans": 2,
                "max_mtu": 9000,
                "vlans": 128
            },
            "usb_jack": "/usb1/1-1/",
            "hwinfo": {
                "esim": false,
                "modem_reset_quirk": false,
                "mt7981_wifi": false,
                "bacnet": false,
                "access_point": false,
                "boot_part": false,
                "sw_offload": false,
                "hi_storage": false,
                "sd_card": false,
                "downstream_kernel": false,
                "gigabit_port": true,
                "dot1x_client": false,
                "port_link": true,
                "ios": true,
                "baseband": false,
                "io": false,
                "ethernet": true,
                "hw_offload": false,
                "rs232": false,
                "power_control": false,
                "bluetooth": false,
                "64mb_ram": false,
                "multi_device": false,
                "dsa": false,
                "mbus": false,
                "at_sim": false,
                "hnat": false,
                "usb": true,
                "micro_usb": false,
                "gateway": false,
                "poe": false,
                "sfp_port": false,
                "usb_port": false,
                "rs485": false,
                "wifi": false,
                "ntrip": false,
                "hw_nat": true,
                "modbus": false,
                "soft_port_mirror": false,
                "industrial-access-point": false,
                "single_port": false,
                "dual_band_ssid": false,
                "2_5_gigabit_port": false,
                "wps": false,
                "nat_offloading": true,
                "dual_modem": false,
                "console": false,
                "basic_router": false,
                "serial": false,
                "testing_kernel": false,
                "serial_reset_quirk": false,
                "guest_wifi": false,
                "tlt_failsafe_boot": false,
                "mobile": false,
                "dual_sim": false,
                "gps": false,
                "bpoffload": false,
                "vendor_wifi": false,
                "qrtrpipes": false,
                "ncm": false,
                "smp": false,
                "rootfs_part": false,
                "verified_boot": false,
                "port_mirror": false,
                "bt": false,
                "tpm": false,
                "pppmobile": false,
                "rndis": false,
                "high_watchdog_priority": false,
                "ledman_lite": false,
                "multi_tag": true
            }
        }
    }
}*/

type RequestDeviceStatusResponseData struct {
	MnfInfo struct {
		MacEth string `json:"macEth"`
		Name   string `json:"name"`
		BlVer  string `json:"blver"`
		HwVer  string `json:"hwver"`
		Serial string `json:"serial"`
		Mac    string `json:"mac"`
	} `json:"mnfnfo"`
	Static struct {
		FwVersion   string `json:"fw_version"`
		Kernel      string `json:"kernel"`
		System      string `json:"system"`
		DeviceName  string `json:"device_name"`
		Hostname    string `json:"hostname"`
		FwBuildDate string `json:"fw_build_date"`
		Model       string `json:"model"`
	} `json:"static"`
}

type RequestDeviceStatusResponse struct {
	Success bool                            `json:"success"`
	Data    RequestDeviceStatusResponseData `json:"data"`
}
