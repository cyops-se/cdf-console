package provisioning

import (
	"fmt"
	"os"
	"path/filepath"
	"server/types"
	"strings"
)

// ScanAvailableFirmware scans the provisioning folder and returns available firmware versions
func ScanAvailableFirmware(provisioningPath string) ([]types.FirmwareVersion, error) {
	var versions []types.FirmwareVersion

	// Get list of model directories
	models, err := os.ReadDir(provisioningPath)
	if err != nil {
		return nil, err
	}

	for _, modelEntry := range models {
		if !modelEntry.IsDir() {
			continue
		}

		modelName := strings.ToUpper(modelEntry.Name())
		if modelName == "SCRIPTS" || modelName == "TMP" {
			continue
		}

		modelPath := filepath.Join(provisioningPath, modelEntry.Name())

		// Get list of version directories for this model
		versionDirs, err := os.ReadDir(modelPath)
		if err != nil {
			continue
		}

		for _, versionEntry := range versionDirs {
			if !versionEntry.IsDir() {
				continue
			}

			versionPath := filepath.Join(modelPath, versionEntry.Name())

			// Find firmware .bin file
			files, err := os.ReadDir(versionPath)
			if err != nil {
				continue
			}

			var firmwareFile, vxlanFile string
			for _, file := range files {
				if strings.HasSuffix(file.Name(), "_WEBUI.bin") {
					firmwareFile = file.Name()
				} else if file.Name() == "vxlan.tar.gz" {
					vxlanFile = file.Name()
				}
			}

			// Only add if firmware file exists
			if firmwareFile != "" {
				versions = append(versions, types.FirmwareVersion{
					Model:        modelName,
					Version:      versionEntry.Name(),
					FirmwareFile: firmwareFile,
					VxlanFile:    vxlanFile,
					Path:         versionPath,
				})
			}
		}
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no firmware for provisioning can be found, check the './provisioning' subfolder to the server!")
	}

	return versions, nil
}

// GetFirmwareVersion finds a specific firmware version by model and version
func GetFirmwareVersion(provisioningPath, model, version string) *types.FirmwareVersion {
	versions, err := ScanAvailableFirmware(provisioningPath)
	if err != nil {
		return nil
	}

	model = strings.ToUpper(model)
	for _, v := range versions {
		if v.Model == model && v.Version == version {
			return &v
		}
	}

	return nil
}
