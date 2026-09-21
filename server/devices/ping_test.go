package devices

import (
	"testing"
	"time"
)

func TestIsValidIPv4(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"127.0.0.1", true},
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"", false},
		{"   ", false},
		{"invalid", false},
		{"256.1.1.1", false},
		{"192.168.1", false},
		{"192.168.1.1.1", false},
		{"::1", false},                     // IPv6
		{"2001:db8::1", false},             // IPv6
		{"fe80::1", false},                 // IPv6
		{"192.168.1.1:8080", false},        // with port
		{"http://192.168.1.1", false},      // with protocol
	}

	for _, test := range tests {
		result := isValidIPv4(test.ip)
		if result != test.expected {
			t.Errorf("isValidIPv4(%q) = %v, expected %v", test.ip, result, test.expected)
		}
	}
}

func TestPingLocalhost(t *testing.T) {
	// Test pinging localhost - should always succeed
	result := Ping("127.0.0.1", QuickPingConfig())

	if !result.Success {
		t.Errorf("Ping to localhost failed: %v", result.Error)
	}

	if result.PacketsRecv == 0 {
		t.Error("No packets received from localhost")
	}

	t.Logf("Ping successful: Sent=%d, Recv=%d, Loss=%.2f%%, RTT=%v",
		result.PacketsSent, result.PacketsRecv, result.PacketLoss, result.RTT)
}

func TestPingInvalidIPv4(t *testing.T) {
	// Test pinging invalid IPv4 addresses
	invalidIPs := []string{
		"",
		"invalid",
		"256.1.1.1",
		"192.168.1",
		"::1",
		"2001:db8::1",
	}

	for _, ip := range invalidIPs {
		result := Ping(ip, QuickPingConfig())
		if result.Success {
			t.Errorf("Ping to invalid IP %q succeeded (should have failed)", ip)
		}
		if result.Error == nil {
			t.Errorf("Ping to invalid IP %q did not return an error", ip)
		}
		t.Logf("Ping to invalid IP %q correctly failed: %v", ip, result.Error)
	}
}

func TestPingNonExistent(t *testing.T) {
	// Test pinging a non-existent IP - should fail or timeout
	config := PingConfig{
		Count:    1,
		Timeout:  100 * time.Millisecond,
		Interval: 0,
		Size:     32,
	}

	result := Ping("192.0.2.1", config) // RFC 5737 TEST-NET-1 (should not respond)

	if result.Success {
		t.Log("Ping succeeded (unexpected but not necessarily an error)")
	} else {
		t.Logf("Ping failed as expected: PacketLoss=%.2f%%", result.PacketLoss)
	}
}

func TestPingOnce(t *testing.T) {
	// Test quick ping to localhost
	success := PingOnce("127.0.0.1", 1*time.Second)

	if !success {
		t.Error("PingOnce to localhost failed")
	}

	t.Log("PingOnce to localhost succeeded")
}

func TestPingOnceInvalid(t *testing.T) {
	// Test PingOnce with invalid IPv4
	success := PingOnce("::1", 1*time.Second)

	if success {
		t.Error("PingOnce with IPv6 address succeeded (should have failed)")
	}

	t.Log("PingOnce with IPv6 address correctly failed")
}

func TestCheckIPPresentICMP(t *testing.T) {
	// Test the drop-in replacement function
	present := CheckIPPresentICMP("127.0.0.1")

	if !present {
		t.Error("CheckIPPresentICMP failed for localhost")
	}

	t.Log("CheckIPPresentICMP for localhost succeeded")
}

func TestCheckIPPresentICMPInvalid(t *testing.T) {
	// Test with invalid addresses
	invalidIPs := []string{
		"",
		"invalid",
		"::1",
		"2001:db8::1",
		"256.1.1.1",
	}

	for _, ip := range invalidIPs {
		present := CheckIPPresentICMP(ip)
		if present {
			t.Errorf("CheckIPPresentICMP(%q) returned true (should be false)", ip)
		}
	}

	t.Log("CheckIPPresentICMP correctly rejected all invalid addresses")
}

func TestCheckIPPresentICMPEmpty(t *testing.T) {
	// Test with empty IP (common case when device IPs are not set)
	present := CheckIPPresentICMP("")

	if present {
		t.Error("CheckIPPresentICMP with empty IP returned true (should be false)")
	}

	t.Log("CheckIPPresentICMP with empty IP correctly returned false")
}

func BenchmarkPing(b *testing.B) {
	config := QuickPingConfig()

	for i := 0; i < b.N; i++ {
		Ping("127.0.0.1", config)
	}
}

func BenchmarkPingOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PingOnce("127.0.0.1", 500*time.Millisecond)
	}
}

func BenchmarkIsValidIPv4(b *testing.B) {
	testIPs := []string{
		"127.0.0.1",
		"192.168.1.1",
		"invalid",
		"::1",
		"",
	}

	for i := 0; i < b.N; i++ {
		for _, ip := range testIPs {
			isValidIPv4(ip)
		}
	}
}
