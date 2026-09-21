package devices

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"server/logger"
	"server/messages"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds configuration for SSH connections
type SSHConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	PrivateKeyPath  string
	Timeout         time.Duration
	HostKeyCallback ssh.HostKeyCallback
}

// NewSSHConfig creates a new SSH configuration with default values
func NewSSHConfig(user, host string) *SSHConfig {
	return &SSHConfig{
		Host:            host,
		Port:            "22",
		User:            user,
		Timeout:         500 * time.Millisecond,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Equivalent to StrictHostKeyChecking=no
	}
}

// getSSHClientConfig creates an SSH client configuration
func (cfg *SSHConfig) getSSHClientConfig() (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	// Try password authentication first if provided
	if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	// Try private key authentication
	if cfg.PrivateKeyPath != "" {
		key, err := os.ReadFile(cfg.PrivateKeyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	// Try default SSH key locations if no specific key was provided or found
	if cfg.PrivateKeyPath == "" || len(authMethods) == 0 {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			defaultKeys := []string{
				filepath.Join(homeDir, ".ssh", "id_rsa"),
				filepath.Join(homeDir, ".ssh", "id_ed25519"),
				filepath.Join(homeDir, ".ssh", "id_ecdsa"),
			}

			for _, keyPath := range defaultKeys {
				key, err := os.ReadFile(keyPath)
				if err == nil {
					signer, err := ssh.ParsePrivateKey(key)
					if err == nil {
						authMethods = append(authMethods, ssh.PublicKeys(signer))
						break
					}
				}
			}
		}
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method available")
	}

	return &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authMethods,
		HostKeyCallback: cfg.HostKeyCallback,
		Timeout:         cfg.Timeout,
	}, nil
}

// connect establishes an SSH connection
func (cfg *SSHConfig) connect() (*ssh.Client, error) {
	clientConfig, err := cfg.getSSHClientConfig()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// RunSSHCommand executes a command over SSH and returns the result
// If useSessionPool is true, it will try to reuse persistent connections
func RunSSHCommand(cfg *SSHConfig, command string) messages.Command {
	return RunSSHCommandWithPool(cfg, command, true)
}

// RunSSHCommandWithPool executes a command over SSH with optional session pooling
func RunSSHCommandWithPool(cfg *SSHConfig, command string, usePool bool) messages.Command {
	logger.Trace("RunSSHCommand", "running command: ssh %s@%s %s ...", cfg.User, cfg.Host, command)

	if usePool {
		// Try to use the session pool
		sm := GetSessionManager()
		session, err := sm.GetSession(cfg.Host, cfg.User)
		if err != nil {
			logger.Error("RunSSHCommand", "failed to get session from pool: %v, falling back to direct connection", err)
			// Fall back to direct connection
			usePool = false
		} else {
			// Use the pooled session
			output, err := session.RunCommand(command)
			result := messages.Command{
				Command: fmt.Sprintf("ssh %s@%s %s", cfg.User, cfg.Host, command),
				Output:  output,
			}
			if err != nil {
				result.Error = err.Error()
				logger.Error("RunSSHCommand", "failed to run SSH command via pool: %s, err: %v", result.Command, result.Error)
			} else {
				logger.Trace("RunSSHCommand", "ssh command completed successfully via pool")
			}
			return result
		}
	}

	// Direct connection (non-pooled)
	client, err := cfg.connect()
	if err != nil {
		logger.Error("RunSSHCommand", "connection failed: %v", err)
		return messages.Command{
			Command: fmt.Sprintf("ssh %s@%s %s", cfg.User, cfg.Host, command),
			Error:   err.Error(),
		}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logger.Error("RunSSHCommand", "failed to create SSH session: %s", err.Error())
		return messages.Command{
			Command: fmt.Sprintf("ssh %s@%s %s", cfg.User, cfg.Host, command),
			Error:   err.Error(),
		}
	}
	defer session.Close()

	var outBuf, errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf

	err = session.Run(command + " || [ $? = 1 ]")
	output := outBuf.String()
	if errBuf.Len() > 0 {
		output += "\n" + errBuf.String()
	}

	result := messages.Command{
		Command: fmt.Sprintf("ssh %s@%s %s", cfg.User, cfg.Host, command),
		Output:  output,
	}

	if err != nil {
		result.Error = err.Error()
		logger.Error("RunSSHCommand", "failed to run SSH command: %s, err: %v, output: %s", result.Command, result.Error, result.Output)
	} else {
		logger.Trace("RunSSHCommand", "ssh command completed successfully")
	}

	return result
}

// SCPUpload uploads a file to a remote host via SCP
func SCPUpload(cfg *SSHConfig, localPath, remotePath string) messages.Command {
	return SCPUploadWithPool(cfg, localPath, remotePath, true)
}

// SCPUploadWithPool uploads a file to a remote host via SCP with optional session pooling
func SCPUploadWithPool(cfg *SSHConfig, localPath, remotePath string, usePool bool) messages.Command {
	logger.Trace("SCPUpload", "uploading %s to %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath)

	// Read the local file
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
			Error:   err.Error(),
		}
	}

	fileData, err := os.ReadFile(localPath)
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
			Error:   err.Error(),
		}
	}

	var client *ssh.Client
	var shouldClose bool

	if usePool {
		// Try to use the session pool
		sm := GetSessionManager()
		session, err := sm.GetSession(cfg.Host, cfg.User)
		if err != nil {
			logger.Error("SCPUpload", "failed to get session from pool: %v, falling back to direct connection", err)
			usePool = false
		} else {
			session.mu.RLock()
			client = session.Client
			session.mu.RUnlock()
			shouldClose = false // Don't close pooled connections
		}
	}

	if !usePool {
		// Direct connection
		var err error
		client, err = cfg.connect()
		if err != nil {
			logger.Error("SCPUpload", "failed to connect via SSH for SCP: %s", err.Error())
			return messages.Command{
				Command: fmt.Sprintf("scp %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
				Error:   err.Error(),
			}
		}
		shouldClose = true
	}

	if shouldClose {
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
			Error:   err.Error(),
		}
	}
	defer session.Close()

	// SCP protocol implementation
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		// Send file header
		fmt.Fprintf(w, "C%04o %d %s\n", fileInfo.Mode().Perm(), len(fileData), filepath.Base(localPath))
		// Send file content
		w.Write(fileData)
		// Send termination byte
		fmt.Fprint(w, "\x00")
	}()

	var outBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &outBuf

	err = session.Run(fmt.Sprintf("scp -t %s", remotePath))
	output := outBuf.String()

	result := messages.Command{
		Command: fmt.Sprintf("scp %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
		Output:  output,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result
}

// SCPDownload downloads a file from a remote host via SCP
func SCPDownload(cfg *SSHConfig, remotePath, localPath string) messages.Command {
	return SCPDownloadWithPool(cfg, remotePath, localPath, true)
}

// SCPDownloadWithPool downloads a file from a remote host via SCP with optional session pooling
func SCPDownloadWithPool(cfg *SSHConfig, remotePath, localPath string, usePool bool) messages.Command {
	logger.Trace("SCPDownload", "downloading %s@%s:%s to %s", cfg.User, cfg.Host, remotePath, localPath)

	var client *ssh.Client
	var shouldClose bool

	if usePool {
		// Try to use the session pool
		sm := GetSessionManager()
		session, err := sm.GetSession(cfg.Host, cfg.User)
		if err != nil {
			logger.Error("SCPDownload", "failed to get session from pool: %v, falling back to direct connection", err)
			usePool = false
		} else {
			session.mu.RLock()
			client = session.Client
			session.mu.RUnlock()
			shouldClose = false
		}
	}

	if !usePool {
		var err error
		client, err = cfg.connect()
		if err != nil {
			logger.Error("SCPDownload", "failed to connect via SSH for SCP: %s", err.Error())
			return messages.Command{
				Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
				Error:   err.Error(),
			}
		}
		shouldClose = true
	}

	if shouldClose {
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
			Error:   err.Error(),
		}
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf

	// Request the file
	w, _ := session.StdinPipe()
	defer w.Close()

	go func() {
		fmt.Fprint(w, "\x00") // Send initial null byte
	}()

	err = session.Run(fmt.Sprintf("scp -f %s", remotePath))
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
			Error:   err.Error(),
		}
	}

	// Parse SCP protocol response
	data := buf.Bytes()
	if len(data) == 0 {
		return messages.Command{
			Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
			Error:   "no data received",
		}
	}

	// Skip the header line (C0644 size filename)
	lines := bytes.SplitN(data, []byte("\n"), 2)
	if len(lines) < 2 {
		return messages.Command{
			Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
			Error:   "invalid SCP response",
		}
	}

	// Extract file content (everything after first newline, excluding trailing null byte)
	fileContent := lines[1]
	if len(fileContent) > 0 && fileContent[len(fileContent)-1] == 0 {
		fileContent = fileContent[:len(fileContent)-1]
	}

	// Ensure the local directory exists
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
			Error:   err.Error(),
		}
	}

	// Write the file
	err = os.WriteFile(localPath, fileContent, 0644)
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
			Error:   err.Error(),
		}
	}

	logger.Trace("SCPDownload", "done")

	return messages.Command{
		Command: fmt.Sprintf("scp %s@%s:%s %s", cfg.User, cfg.Host, remotePath, localPath),
		Output:  fmt.Sprintf("Downloaded %d bytes", len(fileContent)),
	}
}

// SCPUploadDirectory uploads an entire directory to a remote host via SCP
func SCPUploadDirectory(cfg *SSHConfig, localPath, remotePath string) messages.Command {
	return SCPUploadDirectoryWithPool(cfg, localPath, remotePath, true)
}

// SCPUploadDirectoryWithPool uploads an entire directory to a remote host via SCP with optional session pooling
func SCPUploadDirectoryWithPool(cfg *SSHConfig, localPath, remotePath string, usePool bool) messages.Command {
	logger.Trace("SCPUploadDirectory", "uploading directory %s to %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath)

	var client *ssh.Client
	var shouldClose bool

	if usePool {
		// Try to use the session pool
		sm := GetSessionManager()
		session, err := sm.GetSession(cfg.Host, cfg.User)
		if err != nil {
			logger.Error("SCPUploadDirectory", "failed to get session from pool: %v, falling back to direct connection", err)
			usePool = false
		} else {
			session.mu.RLock()
			client = session.Client
			session.mu.RUnlock()
			shouldClose = false
		}
	}

	if !usePool {
		var err error
		client, err = cfg.connect()
		if err != nil {
			logger.Error("SCPUploadDirectory", "failed to connect via SSH for SCP: %s", err.Error())
			return messages.Command{
				Command: fmt.Sprintf("scp -r %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
				Error:   err.Error(),
			}
		}
		shouldClose = true
	}

	if shouldClose {
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("scp -r %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
			Error:   err.Error(),
		}
	}
	defer session.Close()

	var outBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &outBuf

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		filepath.Walk(localPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(localPath, filePath)
			if relPath == "." {
				return nil
			}

			if info.IsDir() {
				fmt.Fprintf(w, "D%04o 0 %s\n", info.Mode().Perm(), filepath.Base(filePath))
			} else {
				fileData, err := os.ReadFile(filePath)
				if err != nil {
					return err
				}
				fmt.Fprintf(w, "C%04o %d %s\n", info.Mode().Perm(), len(fileData), filepath.Base(filePath))
				w.Write(fileData)
				fmt.Fprint(w, "\x00")
			}

			return nil
		})

		fmt.Fprint(w, "E\n") // End directory
	}()

	err = session.Run(fmt.Sprintf("scp -tr %s", remotePath))
	output := outBuf.String()

	logger.Trace("SCPUploadDirectory", "done")

	result := messages.Command{
		Command: fmt.Sprintf("scp -r %s %s@%s:%s", localPath, cfg.User, cfg.Host, remotePath),
		Output:  output,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result
}

// parseTarget parses a target string in the format "user@host" or "user@host:port"
func parseTarget(target string) (user, host, port string) {
	port = "22" // default port

	// Split by @ to get user and host
	parts := strings.SplitN(target, "@", 2)
	if len(parts) == 2 {
		user = parts[0]
		hostPort := parts[1]

		// Check if port is specified
		if strings.Contains(hostPort, ":") {
			hostParts := strings.SplitN(hostPort, ":", 2)
			host = hostParts[0]
			port = hostParts[1]
		} else {
			host = hostPort
		}
	} else {
		host = target
	}

	return
}

// RunSsh executes an SSH command (replacement for the original RunSsh function)
func RunSsh(arg ...string) messages.Command {
	if len(arg) < 2 {
		return messages.Command{
			Command: fmt.Sprintf("ssh %s", strings.Join(arg, " ")),
			Error:   "invalid arguments: expected target and command",
		}
	}

	target := arg[0]
	command := strings.Join(arg[1:], " ")

	user, host, port := parseTarget(target)
	if user == "" {
		user = "root" // default user
	}

	cfg := NewSSHConfig(user, host)
	cfg.Port = port

	return RunSSHCommand(cfg, command)
}

// RunScp executes an SCP file transfer (replacement for the original RunScp function)
func RunScp(arg ...string) messages.Command {
	if len(arg) < 2 {
		return messages.Command{
			Command: fmt.Sprintf("scp %s", strings.Join(arg, " ")),
			Error:   "invalid arguments: expected source and destination",
		}
	}

	source := arg[0]
	dest := arg[1]

	// Determine if it's an upload or download
	if strings.Contains(source, "@") {
		// Download: source is remote, dest is local
		user, host, port := parseTarget(strings.Split(source, ":")[0])
		remotePath := strings.Split(source, ":")[1]

		if user == "" {
			user = "root"
		}

		cfg := NewSSHConfig(user, host)
		cfg.Port = port

		// Handle wildcard downloads
		if strings.Contains(remotePath, "*") {
			return SCPDownload(cfg, remotePath, dest)
		}

		return SCPDownload(cfg, remotePath, dest)
	} else if strings.Contains(dest, "@") {
		// Upload: source is local, dest is remote
		user, host, port := parseTarget(strings.Split(dest, ":")[0])
		remotePath := strings.Split(dest, ":")[1]

		if user == "" {
			user = "root"
		}

		cfg := NewSSHConfig(user, host)
		cfg.Port = port

		// Check if source is a directory
		fileInfo, err := os.Stat(source)
		if err != nil {
			return messages.Command{
				Command: fmt.Sprintf("scp %s %s", source, dest),
				Error:   err.Error(),
			}
		}

		if fileInfo.IsDir() {
			return SCPUploadDirectory(cfg, source, remotePath)
		}

		return SCPUpload(cfg, source, remotePath)
	}

	return messages.Command{
		Command: fmt.Sprintf("scp %s", strings.Join(arg, " ")),
		Error:   "invalid arguments: neither source nor destination is remote",
	}
}

// RunPscp executes an SCP file transfer with password authentication (replacement for the original RunPscp function)
func RunPscp(password string, arg ...string) messages.Command {
	if len(arg) < 2 {
		return messages.Command{
			Command: fmt.Sprintf("pscp %s", strings.Join(arg, " ")),
			Error:   "invalid arguments: expected source and destination",
		}
	}

	source := arg[0]
	dest := arg[1]

	// Parse destination (format: user@host:/path)
	user, host, port := parseTarget(strings.Split(dest, ":")[0])
	remotePath := strings.Split(dest, ":")[1]

	if user == "" {
		user = "root"
	}

	cfg := NewSSHConfig(user, host)
	cfg.Port = port
	cfg.Password = password

	// Check if source is a directory
	fileInfo, err := os.Stat(source)
	if err != nil {
		return messages.Command{
			Command: fmt.Sprintf("pscp %s %s", source, dest),
			Error:   err.Error(),
		}
	}

	if fileInfo.IsDir() {
		return SCPUploadDirectory(cfg, source, remotePath)
	}

	return SCPUpload(cfg, source, remotePath)
}
