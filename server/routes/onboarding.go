package routes

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"server/devices"

	"github.com/gofiber/fiber/v2"
)

func registerOnboardingRoutes(api fiber.Router) {
	// Onboarding routes currently not used by UI
	// Keeping for potential future use
}

func SetupHostSSHKeys() (string, error) {
	hd, _ := os.UserHomeDir()
	sshd := path.Join(hd, ".ssh")
	rsakey := path.Join(sshd, "id_rsa.pub")
	rsaprivatekey := path.Join(sshd, "id_rsa")

	// Note: Using StrictHostKeyChecking=no to avoid known_hosts conflicts
	kh := path.Join(sshd, "known_hosts")
	if _, err := os.Stat(kh); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}

	if _, err := os.Stat(rsakey); errors.Is(err, os.ErrNotExist) {
		devices.RunSystemCommand("ssh-keygen", "-t", "rsa", "-b", "2048", "-N", "", "-f", rsaprivatekey)
	} else if err != nil {
		return "", err
	}

	return rsakey, nil
}


