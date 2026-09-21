package routes

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"server/db"
	"server/devices"
	"server/logger"
	"server/messages"
	"server/provisioning"
	"server/types"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func registerDeviceRoutes(api fiber.Router) {
	api.Get("/device/present/:ip", CheckDeviceIpPresent)
	api.Get("/device/serial/:ip", GetDeviceSerial)
	api.Get("/device/info/:ip", GetDeviceInfo)
	api.Get("/device/firmware/available", GetAvailableFirmware)
	api.Get("/device/firmware/version/:ip", GetFirmwareVersion)
	api.Post("/device/firmware/upgrade", UpgradeFirmware)
	api.Get("/device/vxlan/check/:ip", CheckVXLAN)
	api.Post("/device/vxlan/install", InstallVXLAN)

	api.Post("/device/initssh/:ip", InitDeviceSSHKey)
	api.Post("/device/shell/:ip", RunShellCommand)
	api.Post("/device/script/:ip", RunScript)
	api.Post("/device/script2/:ip", RunScript2)

	api.Get("/device/backup/:id", Backup)
	api.Post("/device/resetstats", ResetAllTrafficStats)

	// SSH session management endpoints
	api.Get("/device/ssh/sessions", GetAllSSHSessions)
	api.Get("/device/ssh/session/:ip", GetSSHSessionStatus)
	api.Delete("/device/ssh/session/:ip", CloseSSHSession)
}

func CheckDeviceIpPresent(c *fiber.Ctx) error {
	ip := c.Params("ip")
	response := messages.RequestResponse{Success: false}
	response.Success = devices.CheckIPPresent(ip)

	return c.Status(http.StatusOK).JSON(response)
}

func InitDeviceSSHKey(c *fiber.Ctx) error {
	ip := c.Params("ip")

	var loginRequest messages.LoginRequest
	var response messages.CommandResponse
	response.Success = true

	if err := json.Unmarshal(c.Body(), &loginRequest); err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error())
		return c.Status(http.StatusOK).JSON(response)
	}

	logger.Log("trace", "Login request", fmt.Sprintf("%+v", loginRequest))

	rsakey, err := SetupHostSSHKeys()
	if err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error())
		return c.Status(http.StatusOK).JSON(response)
	}

	// Copy SSH file to make subsequent request password-less
	command := devices.RunPscp(loginRequest.Password, rsakey, "root@"+ip+":/etc/dropbear/authorized_keys")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	command = devices.RunSsh("root@"+ip, "uci delete vuci.main.firstlogin && uci commit")
	response.Commands = append(response.Commands, command)
	if command.Error != "" && !strings.Contains(command.Output, "Entry not found") {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func RunScript(c *fiber.Ctx) error {
	ip := c.Params("ip")

	var response messages.CommandResponse
	response.Success = true
	var scriptRequest messages.RunScriptRequest

	if err := json.Unmarshal(c.Body(), &scriptRequest); err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error())
		return c.Status(http.StatusOK).JSON(response)
	}

	fullname := "./provisioning/scripts/" + scriptRequest.Script
	logger.Log("trace", "Script fullname", fullname)

	// Copy SSH file to make subsequent request password-less
	command := devices.RunScp(fullname, "root@"+ip+":/tmp/"+scriptRequest.Script)
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	cmd := "sh /tmp/" + scriptRequest.Script + " " + scriptRequest.Arguments
	logger.Log("trace", "Command", cmd)

	command = devices.RunSsh("root@"+ip, cmd)
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func RunScript2(c *fiber.Ctx) error {
	ip := c.Params("ip")

	var response messages.CommandResponse
	response.Success = true
	var scriptRequest messages.RunScript2Request
	if err := json.Unmarshal(c.Body(), &scriptRequest); err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error())
		return c.Status(http.StatusOK).JSON(response)
	}

	fullname := "./provisioning/scripts/" + scriptRequest.Script
	content, err := os.ReadFile(fullname)
	if err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error())
		return c.Status(http.StatusOK).JSON(response)
	}

	strcontent := string(content)
	for _, kv := range scriptRequest.Arguments {
		strcontent = strings.ReplaceAll(strcontent, "["+kv.Key+"]", kv.Value)
	}

	os.MkdirAll("./provisioning/tmp/", os.ModePerm)
	fullname = "./provisioning/tmp/" + scriptRequest.Script
	os.WriteFile(fullname, []byte(strcontent), 0644)

	// Copy SSH file to make subsequent request password-less
	command := devices.RunScp(fullname, "root@"+ip+":/tmp/"+scriptRequest.Script)
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	cmd := "sh /tmp/" + scriptRequest.Script
	logger.Log("trace", "Command", cmd)

	command = devices.RunSsh("root@"+ip, cmd)
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func RunShellCommand(c *fiber.Ctx) error {
	var response messages.CommandResponse
	response.Success = true

	ip := c.Params("ip")
	cmd := string(c.Body())
	logger.Log("trace", "Command", cmd)

	command := devices.RunSsh("root@"+ip, cmd)
	response.Commands = append(response.Commands, command)
	if command.Error != "" && !strings.Contains(command.Output, "Entry not found") {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func GetDeviceSerial(c *fiber.Ctx) error {
	var response messages.CommandResponse
	response.Success = true

	ip := c.Params("ip")

	// Check serial number
	target := fmt.Sprintf("root@%s", ip)
	command := devices.RunSsh(target, "mnf_info -s")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func GetDeviceInfo(c *fiber.Ctx) error {
	ip := c.Params("ip")
	target := fmt.Sprintf("root@%s", ip)

	deviceInfo := types.DeviceInfo{}

	// Get serial number
	serialCmd := devices.RunSsh(target, "mnf_info -s")
	if serialCmd.Error == "" {
		deviceInfo.SerialNo = strings.TrimSpace(serialCmd.Output)
	}

	// Get model
	modelCmd := devices.RunSsh(target, "uci show system.system.device_code")
	if modelCmd.Error == "" {
		parts := strings.Split(strings.ReplaceAll(modelCmd.Output, "'", ""), "=")
		if len(parts) > 1 {
			deviceInfo.Model = strings.TrimSpace(parts[1])[0:4]
		}
	}

	// Get firmware version
	firmwareCmd := devices.RunSsh(target, "uci show system.system.device_fw_version")
	if firmwareCmd.Error == "" {
		parts := strings.Split(strings.ReplaceAll(firmwareCmd.Output, "'", ""), "=")
		if len(parts) > 1 {
			deviceInfo.Firmware = strings.TrimSpace(parts[1])
		}
	}

	// Get hostname
	hostnameCmd := devices.RunSsh(target, "uci show system.system.hostname")
	if hostnameCmd.Error == "" {
		parts := strings.Split(strings.ReplaceAll(hostnameCmd.Output, "'", ""), "=")
		if len(parts) > 1 {
			deviceInfo.Hostname = strings.TrimSpace(parts[1])
		}
	}

	// Get 4G IP address if available
	apn4gCmd := devices.GetDeviceApn4G(ip)
	if apn4gCmd.Error == "" {
		deviceInfo.Apn4GCIDR = strings.TrimSuffix(apn4gCmd.Output, "\n")
	}

	return c.Status(http.StatusOK).JSON(deviceInfo)
}

func GetAvailableFirmware(c *fiber.Ctx) error {
	response := types.FirmwareListResponse{
		Success: true,
	}

	versions, err := provisioning.ScanAvailableFirmware("./provisioning")
	if err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response.Versions = versions
	return c.Status(http.StatusOK).JSON(response)
}

func GetFirmwareVersion(c *fiber.Ctx) error {
	ip := c.Params("ip")
	var response messages.CommandResponse
	response.Success = true

	// Check firmware version
	command := devices.RunSsh("root@"+ip, "uci show system.system.device_fw_version")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	}

	return c.Status(http.StatusOK).JSON(response)
}

type FirmwareUpgradeRequest struct {
	IP      string `json:"ip"`
	Model   string `json:"model"`
	Version string `json:"version"`
}

func UpgradeFirmware(c *fiber.Ctx) error {
	var request FirmwareUpgradeRequest
	var response messages.CommandResponse
	response.Success = true

	if err := c.BodyParser(&request); err != nil {
		response.Success = false
		response.Errors = append(response.Errors, "Invalid request: "+err.Error())
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Get firmware version info
	firmware := provisioning.GetFirmwareVersion("./provisioning", request.Model, request.Version)
	if firmware == nil {
		response.Success = false
		response.Errors = append(response.Errors, fmt.Sprintf("Firmware not found for model %s version %s", request.Model, request.Version))
		return c.Status(http.StatusNoContent).JSON(response)
	}

	src := filepath.Join(firmware.Path, firmware.FirmwareFile)
	cmd := fmt.Sprintf("sysupgrade /tmp/setup/%s", firmware.FirmwareFile)

	// Create setup directory
	command := devices.RunSsh("root@"+request.IP, "mkdir -p /tmp/setup")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	// Copy firmware to device
	command = devices.RunScp(src, "root@"+request.IP+":/tmp/setup")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	// Install firmware
	command = devices.RunSsh("root@"+request.IP, cmd)
	response.Commands = append(response.Commands, command)
	if command.Error != "" && !strings.Contains(command.Error, "exit status 0xffffffff") && !strings.Contains(command.Error, "remote command exited without") {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func CheckVXLAN(c *fiber.Ctx) error {
	ip := c.Params("ip")
	var response messages.CommandResponse
	response.Success = true

	// Check firmware version
	command := devices.RunSsh("root@"+ip, "lsmod | grep -i vxlan")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
	} else {
		if !strings.Contains(command.Output, "vxlan") {
			response.Success = false
			response.Errors = append(response.Errors, "vxlan modules not found")
		}
	}

	return c.Status(http.StatusOK).JSON(response)
}

type VXLANInstallRequest struct {
	IP      string `json:"ip"`
	Model   string `json:"model"`
	Version string `json:"version"`
}

func InstallVXLAN(c *fiber.Ctx) error {
	var request VXLANInstallRequest
	var response messages.CommandResponse
	response.Success = true

	if err := c.BodyParser(&request); err != nil {
		response.Success = false
		response.Errors = append(response.Errors, "Invalid request: "+err.Error())
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Get firmware version info
	firmware := provisioning.GetFirmwareVersion("./provisioning", request.Model, request.Version)
	if firmware == nil {
		response.Success = false
		response.Errors = append(response.Errors, fmt.Sprintf("Firmware not found for model %s version %s", request.Model, request.Version))
		return c.Status(http.StatusNoContent).JSON(response)
	}

	if firmware.VxlanFile == "" {
		response.Success = false
		response.Errors = append(response.Errors, fmt.Sprintf("VXLAN package not available for model %s version %s", request.Model, request.Version))
		return c.Status(http.StatusNoContent).JSON(response)
	}

	src := filepath.Join(firmware.Path, firmware.VxlanFile)

	// Create setup directory
	command := devices.RunSsh("root@"+request.IP, "mkdir -p /tmp/setup")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	// Copy VXLAN installation files
	command = devices.RunScp(src, "root@"+request.IP+":/tmp/setup/")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	// Unpack VXLAN installation files
	command = devices.RunSsh("root@"+request.IP, "cd /tmp/setup && tar -xzf vxlan.tar.gz")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	// Install VXLAN installation files
	command = devices.RunSsh("root@"+request.IP, "cd /tmp/setup && opkg --conf /etc/opkg.conf --force-removal-of-essential-packages --tmp-dir /tmp/setup install *.ipk && chown root.root /etc")
	response.Commands = append(response.Commands, command)
	if command.Error != "" {
		response.Success = false
		response.Errors = append(response.Errors, command.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func backupDevice(id string) error {
	var device types.Device
	result := db.DB.Preload(clause.Associations).First(&device, id)
	if result.Error != nil {
		return result.Error
	}

	if device.Role == "hub" {
		for _, d := range device.Endpoints {
			backupDevice(fmt.Sprintf("%d", d.ID))
		}
	}

	logger.Log("trace", "Backing up device", fmt.Sprintf("%s [%s]", device.Wan, device.Role))
	devices.RunSsh("root@"+device.PreferredIp, "rm /tmp/backup-*")
	command := devices.RunSsh("root@"+device.PreferredIp, "sysupgrade -b /tmp/backup-${HOSTNAME}-$(date +%F).tar.gz")
	if command.Error != "" {
		return fmt.Errorf(command.Error)
	}

	os.MkdirAll("./backup/", os.ModePerm)
	command = devices.RunScp("root@"+device.PreferredIp+":/tmp/backup-*", "./backup/")
	if command.Error != "" {
		return fmt.Errorf(command.Error)
	}

	return nil
}

func Backup(c *fiber.Ctx) error {
	id := c.Params("id")
	var response messages.CommandResponse
	response.Success = true

	if err := backupDevice(id); err != nil {
		response.Success = false
		response.Errors = append(response.Errors, err.Error)
		return c.Status(http.StatusOK).JSON(response)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func calcMaskSize(mask string) int {
	var b1, b2, b3, b4 int
	mask = strings.ReplaceAll(mask, ".", " ")
	if _, err := fmt.Sscanf(mask, "'%d %d %d %d'", &b1, &b2, &b3, &b4); err != nil {
		logger.Log("error", "Failed to parse mask", err.Error())
	}

	logger.Log("trace", "Mask parsed", fmt.Sprintf("%s = %d.%d.%d.%d", mask, b1, b2, b3, b4))
	sz, _ := net.IPv4Mask(byte(b1), byte(b2), byte(b3), byte(b4)).Size()
	return sz
}

func parseUciOutput(ucioutput string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(ucioutput, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) > 1 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func ResetAllTrafficStats(c *fiber.Ctx) error {
	var response messages.RequestResponse
	response.Success = true

	// Reset all traffic stats for all devices
	result := db.DB.Model(&types.Device{}).Updates(map[string]interface{}{
		"rx_bytes":   0,
		"tx_bytes":   0,
		"rx_packets": 0,
		"tx_packets": 0,
	})

	if result.Error != nil {
		response.Success = false
		response.Errors = append(response.Errors, result.Error.Error())
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	logger.Log("trace", "Reset traffic stats", fmt.Sprintf("Reset stats for %d devices", result.RowsAffected))
	return c.Status(http.StatusOK).JSON(response)
}

// GetAllSSHSessions returns information about all active SSH sessions
func GetAllSSHSessions(c *fiber.Ctx) error {
	sm := devices.GetSessionManager()
	sessions := sm.GetAllSessionsInfo()

	type Response struct {
		Success  bool                           `json:"success"`
		Sessions map[string]devices.SessionInfo `json:"sessions"`
	}

	response := Response{
		Success:  true,
		Sessions: sessions,
	}

	return c.Status(http.StatusOK).JSON(response)
}

// GetSSHSessionStatus returns the status of a specific SSH session
func GetSSHSessionStatus(c *fiber.Ctx) error {
	ip := c.Params("ip")
	sm := devices.GetSessionManager()
	status := sm.GetSessionStatus(ip)

	type Response struct {
		Success bool                  `json:"success"`
		IP      string                `json:"ip"`
		Status  devices.SessionStatus `json:"status"`
	}

	response := Response{
		Success: true,
		IP:      ip,
		Status:  status,
	}

	return c.Status(http.StatusOK).JSON(response)
}

// CloseSSHSession closes an SSH session for a specific IP
func CloseSSHSession(c *fiber.Ctx) error {
	ip := c.Params("ip")
	sm := devices.GetSessionManager()
	sm.RemoveSession(ip)

	var response messages.RequestResponse
	response.Success = true

	logger.Log("info", "SSH session closed", fmt.Sprintf("Closed SSH session for %s", ip))
	return c.Status(http.StatusOK).JSON(response)
}
