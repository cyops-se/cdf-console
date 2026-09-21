package devices

import (
	"fmt"
	"server/db"
	"server/logger"
	"server/messages"
	"server/types"
	"strconv"
	"strings"
	"time"
)

var Devices map[string]types.Device
var DefaultAddressPresent bool

func CheckDefaultPresent() bool {
	// Check if the default device at 192.168.1.1 is present using ICMP ping
	return CheckIPPresent("192.168.1.1")
}

func CheckIPPresent(ip string) bool {
	// Use ICMP ping instead of HTTP GET for presence detection
	// This is more reliable and works regardless of HTTP service status
	return CheckIPPresentICMP(ip)
}

func GetDeviceSerial(ip string) messages.Command {
	// Check serial number
	target := fmt.Sprintf("root@%s", ip)
	command := RunSsh(target, "mnf_info -s")
	return command
}

func GetDeviceApn4G(ip string) messages.Command {
	target := fmt.Sprintf("root@%s", ip)
	command := RunSsh(target, "ip a | grep qmimux | awk '{print $2}' | grep -v qmi")
	return command
}

func CheckHostname(ip string, hostname string) messages.Command {
	// Check serial number
	target := fmt.Sprintf("root@%s", ip)
	command := RunSsh(target, "echo ${HOSTNAME}")
	command.Output = strings.ReplaceAll(command.Output, "\n", "")
	if command.Output != hostname {
		logger.Trace("CheckHostname", "setting hostname from %s to %s", command.Output, hostname)
		command = RunSsh("root@"+ip, fmt.Sprintf("uci set system.system.hostname='%s' && uci set system.system.devicename='%s' && uci commit && /etc/init.d/system reload", hostname, hostname))
	}
	return command
}

// CollectTrafficStats collects network traffic statistics for all operational interfaces on the device
// Returns a map of interface name -> stats
func CollectTrafficStats(ip string) (map[string]types.InterfaceStats, bool) {
	target := fmt.Sprintf("root@%s", ip)
	// Get statistics, operational state, and IP address for all network interfaces
	// Only collect data for interfaces that are operational (operstate == "up")
	command := RunSsh(target, "for iface in /sys/class/net/*; do if [ -d \"$iface/statistics\" ]; then name=$(basename $iface); operstate=$(cat \"$iface/operstate\" 2>/dev/null || echo unknown); if [ \"$operstate\" = \"up\" ]; then echo \"$name\"; echo \"$operstate\"; ipaddr=$(ip -4 addr show \"$name\" 2>/dev/null | awk '/inet / {split($2,a,\"/\"); print a[1]; exit}'); echo \"${ipaddr:-none}\"; cat \"$iface/statistics/rx_bytes\" \"$iface/statistics/tx_bytes\" \"$iface/statistics/rx_packets\" \"$iface/statistics/tx_packets\"; fi; fi; done")

	if command.Error != "" {
		logger.Error("CollectTrafficStats", "failed to collect traffic stats from %s: %s", ip, command.Error)
		return nil, false
	}

	// Parse the output
	lines := strings.Split(strings.TrimSpace(command.Output), "\n")
	interfaceStats := make(map[string]types.InterfaceStats)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check if this is an interface name (non-numeric)
		if _, err := strconv.ParseUint(line, 10, 64); err != nil {
			// This is an interface name, next 6 lines should be: operstate, ipaddress, rx_bytes, tx_bytes, rx_packets, tx_packets
			if i+6 < len(lines) {
				ifaceName := line
				operstate := strings.TrimSpace(lines[i+1])

				// Skip loopback interface
				if ifaceName == "lo" {
					i += 6
					continue
				}

				var stat types.InterfaceStats
				stat.Operstate = operstate

				if rxBytesVal, err := strconv.ParseUint(strings.TrimSpace(lines[i+3]), 10, 64); err == nil {
					stat.RxBytes = rxBytesVal
				}
				if txBytesVal, err := strconv.ParseUint(strings.TrimSpace(lines[i+4]), 10, 64); err == nil {
					stat.TxBytes = txBytesVal
				}
				if rxPacketsVal, err := strconv.ParseUint(strings.TrimSpace(lines[i+5]), 10, 64); err == nil {
					stat.RxPackets = rxPacketsVal
				}
				if txPacketsVal, err := strconv.ParseUint(strings.TrimSpace(lines[i+6]), 10, 64); err == nil {
					stat.TxPackets = txPacketsVal
				}

				interfaceStats[ifaceName] = stat
				i += 6 // Skip the lines we just processed
			}
		}
	}

	logger.Trace("CollectTrafficStats", "collected stats from %s for %d operational interfaces", ip, len(interfaceStats))

	return interfaceStats, true
}

func CheckWgStatusOnHub(hub types.Device) {
	target := fmt.Sprintf("root@%s", hub.PreferredIp)
	command := RunSsh(target, "wg")
	lines := strings.Split(command.Output, "\n")
	wgactive := false

	for i, l := range lines {
		if strings.HasPrefix(l, "peer:") {
			var device types.Device
			for goon := true; goon; i++ {
				if lines[i] == "" {
					break
				} else if strings.HasPrefix(lines[i], "  endpoint:") {
					ip := strings.TrimSpace(strings.Split(lines[i], ":")[1])
					db.DB.Where("wan == ?", ip).Or("apn4_g == ?", ip).Or("lan == ?", ip).First(&device)
					if device.ID > 0 {
						device.WgStatus = "not active"
					}
				} else if strings.HasPrefix(lines[i], "  latest handshake:") {
					if device.ID > 0 {
						device.WgStatus = "active"
						wgactive = true
					}
				}
			}

			if device.ID > 0 {
				db.DB.Model(&device).Updates(&types.Device{WgStatus: device.WgStatus})
			}
		}
	}

	hub.WgStatus = "not active"
	if wgactive {
		hub.WgStatus = "active"
	}

	db.DB.Model(&hub).Updates(&types.Device{WgStatus: hub.WgStatus})
}

func CheckPresenceOfKnown() {
	var devices []types.Device
	if result := db.DB.Find(&devices); result.Error != nil {
		logger.Error("CheckPresenceOfKnown", "unable to get devices from database")
		return
	}

	for _, d := range devices {
		d.Status = "unknown"
		present := false

		// Check WAN interface
		d.WanPresent = CheckIPPresent(d.Wan)
		if d.WanPresent {
			d.PreferredIp = d.Wan
			// Capture MAC address from ARP table after successful ping
			// d.WanInterface.MacAddress = GetMacAddressFromARP(d.Wan)
		}

		// Check APN interface
		d.Apn4GPresent = CheckIPPresent(d.Apn4G)
		if d.Apn4GPresent && !present {
			d.PreferredIp = d.Apn4G
		}
		if d.Apn4GPresent {
			// Capture MAC address from ARP table after successful ping
			// d.Apn4GInterface.MacAddress = GetMacAddressFromARP(d.Apn4G)
		}

		// Check LAN interface
		d.LanPresent = CheckIPPresent(d.Lan)
		if d.LanPresent && !present {
			d.PreferredIp = d.Lan
		}
		if d.LanPresent {
			// Capture MAC address from ARP table after successful ping
			// d.LanInterface.MacAddress = GetMacAddressFromARP(d.Lan)
		}

		if d.Apn4GPresent || d.WanPresent || d.LanPresent {
			d.Status = "reachable"
			present = true
			GetDeviceSerial(d.PreferredIp)
			d.SshStatus = string(GetSessionManager().GetSessionStatus(d.PreferredIp))
			logger.Trace("ssh check", "ssh session for ip %s has status %s", d.PreferredIp, d.SshStatus)
		} else {
			d.Status = "not reachable"
			d.WgStatus = "unknown"
			d.SshStatus = "unknown"
		}

		d.LastCheck = time.Now().UTC()

		// Have to use a map[string]interface{} to update zero/empty fields :(
		db.DB.Model(&d).Updates(map[string]interface{}{"wan_present": d.WanPresent, "lan_present": d.LanPresent, "apn4_g_present": d.Apn4GPresent,
			"status": d.Status, "ssh_status": d.SshStatus, "preferred_ip": d.PreferredIp, "last_check": d.LastCheck, "wg_status": d.WgStatus})

		if present && d.Role == "hub" {
			CheckWgStatusOnHub(d)
		}
	}
}
