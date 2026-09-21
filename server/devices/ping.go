package devices

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"server/logger"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// PingResult contains the result of a ping operation
type PingResult struct {
	Success     bool
	RTT         time.Duration
	PacketsSent int
	PacketsRecv int
	PacketLoss  float64
	Error       error
}

// PingConfig holds configuration for ping operations
type PingConfig struct {
	Count    int           // Number of packets to send (default: 3)
	Timeout  time.Duration // Timeout per packet (default: 1 second)
	Interval time.Duration // Interval between packets (default: 100ms)
	Size     int           // Size of ping payload (default: 56 bytes)
}

// DefaultPingConfig returns a ping configuration with sensible defaults
func DefaultPingConfig() PingConfig {
	return PingConfig{
		Count:    3,
		Timeout:  1 * time.Second,
		Interval: 100 * time.Millisecond,
		Size:     56,
	}
}

// QuickPingConfig returns a fast ping configuration for presence checking
func QuickPingConfig() PingConfig {
	return PingConfig{
		Count:    2,
		Timeout:  50 * time.Millisecond,
		Interval: 50 * time.Millisecond,
		Size:     32,
	}
}

// isValidIPv4 checks if the given string is a valid IPv4 address
func isValidIPv4(ip string) bool {
	// Check if it's empty or contains invalid characters
	if ip == "" || strings.TrimSpace(ip) == "" {
		return false
	}

	// Parse the IP
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check if it's an IPv4 address (not IPv6)
	return parsedIP.To4() != nil
}

// Ping sends ICMP echo requests to the specified IPv4 host
// This implementation is compliant with RFC 792 (Internet Control Message Protocol)
// https://tools.ietf.org/html/rfc792
func Ping(host string, config PingConfig) PingResult {
	result := PingResult{
		PacketsSent: config.Count,
	}

	// Validate that the host is a valid IPv4 address
	if !isValidIPv4(host) {
		result.Error = fmt.Errorf("invalid IPv4 address: %s", host)
		return result
	}

	// Resolve the host to an IP address
	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		result.Error = fmt.Errorf("failed to resolve IPv4 host: %w", err)
		return result
	}

	// Verify it's IPv4
	if ipAddr.IP.To4() == nil {
		result.Error = fmt.Errorf("not an IPv4 address: %s", host)
		return result
	}

	// Open ICMP connection for IPv4
	conn, err := icmp.ListenPacket("ip4:icmp", "")
	if err != nil {
		result.Error = fmt.Errorf("failed to listen for ICMP (may need elevated privileges): %w", err)
		return result
	}
	defer conn.Close()

	// Generate a random identifier for this ping session
	pid := os.Getpid() & 0xffff

	var totalRTT time.Duration
	for seq := 1; seq <= config.Count; seq++ {
		// Create ICMP echo request message
		// RFC 792 specifies:
		// - Type: 8 (Echo Request)
		// - Code: 0
		// - Checksum: calculated
		// - Identifier: process ID
		// - Sequence Number: incremental

		// Create payload with random data
		payload := make([]byte, config.Size)
		rand.Read(payload)

		// Add timestamp to payload for RTT calculation
		if len(payload) >= 8 {
			binary.BigEndian.PutUint64(payload[:8], uint64(time.Now().UnixNano()))
		}

		// Build ICMP message
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho, // Type 8: Echo Request
			Code: 0,
			Body: &icmp.Echo{
				ID:   pid,
				Seq:  seq,
				Data: payload,
			},
		}

		// Marshal the message to bytes (this calculates the checksum)
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			logger.Error("Ping", "failed to marshal ICMP message: %v", err)
			continue
		}

		// Send the ICMP echo request
		start := time.Now()
		_, err = conn.WriteTo(msgBytes, ipAddr)
		if err != nil {
			logger.Error("Ping", "failed to send ICMP packet: %v", err)
			continue
		}

		// Set read deadline for timeout
		conn.SetReadDeadline(time.Now().Add(config.Timeout))

		// Wait for echo reply
		reply := make([]byte, 1500) // MTU size
		for {
			n, peer, err := conn.ReadFrom(reply)
			if err != nil {
				// Timeout or other error
				break
			}

			// Parse the ICMP reply
			replyMsg, err := icmp.ParseMessage(1, reply[:n]) // Protocol 1 = ICMP
			if err != nil {
				continue
			}

			// Check if this is an echo reply for our request
			// RFC 792: Echo Reply has Type 0, Code 0
			echoReply, ok := replyMsg.Body.(*icmp.Echo)
			if !ok {
				continue
			}

			// Verify this reply is for our ping
			if echoReply.ID != pid || echoReply.Seq != seq {
				continue
			}

			// Verify the reply is from the host we pinged
			if peer.String() != ipAddr.String() {
				continue
			}

			// Check message type is Echo Reply (Type 0)
			if replyMsg.Type != ipv4.ICMPTypeEchoReply {
				continue
			}

			// Calculate RTT
			rtt := time.Since(start)
			totalRTT += rtt
			result.PacketsRecv++
			result.Success = true

			break
		}

		// Wait before sending next packet (if not the last one)
		if seq < config.Count {
			time.Sleep(config.Interval)
		}
	}

	// Calculate statistics
	if result.PacketsRecv > 0 {
		result.RTT = totalRTT / time.Duration(result.PacketsRecv)
	}

	result.PacketLoss = float64(result.PacketsSent-result.PacketsRecv) / float64(result.PacketsSent) * 100.0

	return result
}

// PingOnce sends a single ICMP echo request and returns quickly
// This is optimized for fast presence detection
// Only pings valid IPv4 addresses
func PingOnce(host string, timeout time.Duration) bool {
	// Validate IPv4 before attempting to ping
	if !isValidIPv4(host) {
		return false
	}

	config := PingConfig{
		Count:    1,
		Timeout:  timeout,
		Interval: 0,
		Size:     32,
	}

	result := Ping(host, config)
	return result.Success
}

// CheckIPPresentICMP checks if an IP address is reachable using ICMP ping
// This is a drop-in replacement for CheckIPPresent that uses ICMP instead of HTTP
// Only pings valid IPv4 addresses
func CheckIPPresentICMP(ip string) bool {
	// Validate IPv4 before attempting to ping
	if !isValidIPv4(ip) {
		return false
	}

	// Quick check with 2 pings, 500ms timeout each
	config := QuickPingConfig()
	result := Ping(ip, config)

	// Consider it present if we got at least one reply
	return result.Success
}

// GetMacAddressFromARP retrieves the MAC address for an IP from the ARP table
// This function should be called after a successful ping to ensure the ARP entry exists
// Returns the MAC address in format "aa:bb:cc:dd:ee:ff" or empty string if not found
func GetMacAddressFromARP(ip string) string {
	// Validate IPv4
	if !isValidIPv4(ip) {
		return ""
	}

	var cmd *exec.Cmd
	var macRegex *regexp.Regexp

	// Different ARP command and output format based on OS
	switch runtime.GOOS {
	case "windows":
		// Windows: arp -a
		// Example output: 192.168.1.1       aa-bb-cc-dd-ee-ff     dynamic
		cmd = exec.Command("arp", "-a", ip)
		macRegex = regexp.MustCompile(`([0-9a-fA-F]{2}[-:]){5}[0-9a-fA-F]{2}`)
	case "linux", "darwin":
		// Linux/Mac: arp -n
		// Example output: 192.168.1.1 ether aa:bb:cc:dd:ee:ff C eth0
		cmd = exec.Command("arp", "-n", ip)
		macRegex = regexp.MustCompile(`([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}`)
	default:
		logger.Trace("GetMacAddressFromARP", "unsupported OS: %s", runtime.GOOS)
		return ""
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Trace("GetMacAddressFromARP", "failed to run arp command for %s: %v", ip, err)
		return ""
	}

	// Search for MAC address in output
	mac := macRegex.FindString(string(output))
	if mac == "" {
		logger.Trace("GetMacAddressFromARP", "no MAC address found in ARP table for %s", ip)
		return ""
	}

	// Normalize MAC address format to use colons (aa:bb:cc:dd:ee:ff)
	mac = strings.ToLower(mac)
	mac = strings.ReplaceAll(mac, "-", ":")

	logger.Trace("GetMacAddressFromARP", "found MAC address %s for IP %s", mac, ip)
	return mac
}
