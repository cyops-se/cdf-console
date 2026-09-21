package types

// FirmwareVersion represents a firmware version available in the provisioning folder
type FirmwareVersion struct {
	Model        string `json:"model"`        // e.g., "RUTX", "RUTM", "RUT2M", "RUT9M"
	Version      string `json:"version"`      // e.g., "7.12", "7.17.5"
	FirmwareFile string `json:"firmwarefile"` // Filename of the .bin file
	VxlanFile    string `json:"vxlanfile"`    // Filename of vxlan.tar.gz
	Path         string `json:"path"`         // Full path to the firmware folder
}

// DeviceInfo represents detailed information about a device
type DeviceInfo struct {
	SerialNo  string `json:"serialno"`
	Model     string `json:"model"`
	Firmware  string `json:"firmware"`
	Hostname  string `json:"hostname"`
	Apn4GCIDR string `json:"apn4gcidr"`
}

// FirmwareListResponse contains available firmware versions
type FirmwareListResponse struct {
	Success  bool              `json:"success"`
	Versions []FirmwareVersion `json:"versions"`
	Errors   []string          `json:"errors"`
}
