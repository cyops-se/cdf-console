package devices

import "time"

func CheckDefaultPresentWorker() {
	// Check if 192.168.1.1 is present every 5 secs
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		DefaultAddressPresent = CheckDefaultPresent()
	}
}

func CheckAvailabilityWorker() {
	ticker := time.NewTicker(15 * time.Second)

	for range ticker.C {
		CheckPresenceOfKnown()
	}
}
