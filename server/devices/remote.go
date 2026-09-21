package devices

import (
	"os/exec"
	"server/logger"
)

func RunSystemCommand(name string, arg ...string) string {
	prc := exec.Command(name, arg...)

	out, err := prc.CombinedOutput()
	if err != nil {
		logger.Error("RunSystemCommand", "system command failed: %s", err.Error())
	}

	logger.Trace("RunSystemCommand", "out: %s, err: %v", string(out), err)
	return string(out)
}

// Old external command implementations removed - now using native Go SSH/SCP client in ssh_client.go
// RunPscp, RunScp, and RunSsh are now implemented in ssh_client.go using golang.org/x/crypto/ssh
