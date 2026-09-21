package devices

import (
	"fmt"
	"server/logger"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// SessionStatus represents the health status of an SSH session
type SessionStatus string

const (
	StatusValid      SessionStatus = "valid"       // Session is active and responsive
	StatusInvalid    SessionStatus = "invalid"     // Session failed validation
	StatusDisconnect SessionStatus = "disconnected" // Session is disconnected
	StatusUnknown    SessionStatus = "unknown"     // Session status is not yet determined
)

// SSHSession represents a persistent SSH connection to a device
type SSHSession struct {
	Client       *ssh.Client
	Config       *SSHConfig
	Status       SessionStatus
	LastUsed     time.Time
	LastValidate time.Time
	ConnectedAt  time.Time
	FailCount    int
	mu           sync.RWMutex
}

// SessionManager manages persistent SSH connections for multiple devices
type SessionManager struct {
	sessions           map[string]*SSHSession
	mu                 sync.RWMutex
	validationInterval time.Duration
	staleTimeout       time.Duration
	maxFailCount       int
	stopChan           chan struct{}
	wg                 sync.WaitGroup
}

var (
	// GlobalSessionManager is the singleton instance
	GlobalSessionManager *SessionManager
	sessionManagerOnce   sync.Once
)

// GetSessionManager returns the singleton SessionManager instance
func GetSessionManager() *SessionManager {
	sessionManagerOnce.Do(func() {
		GlobalSessionManager = NewSessionManager()
		GlobalSessionManager.Start()
	})
	return GlobalSessionManager
}

// NewSessionManager creates a new SSH session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:           make(map[string]*SSHSession),
		validationInterval: 30 * time.Second, // Validate sessions every 30 seconds
		staleTimeout:       5 * time.Minute,  // Close sessions unused for 5 minutes
		maxFailCount:       3,                // Max consecutive validation failures before marking invalid
		stopChan:           make(chan struct{}),
	}
}

// Start begins background tasks for session management
func (sm *SessionManager) Start() {
	sm.wg.Add(2)
	go sm.validationWorker()
	go sm.cleanupWorker()
	logger.Log("info", "SSH Session Manager started", "")
}

// Stop gracefully shuts down the session manager
func (sm *SessionManager) Stop() {
	close(sm.stopChan)
	sm.wg.Wait()

	// Close all sessions
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for ip, session := range sm.sessions {
		session.Close()
		delete(sm.sessions, ip)
	}
	logger.Log("info", "SSH Session Manager stopped", "")
}

// GetSession retrieves an existing session or creates a new one
func (sm *SessionManager) GetSession(ip string, user string) (*SSHSession, error) {
	sm.mu.RLock()
	session, exists := sm.sessions[ip]
	sm.mu.RUnlock()

	if exists {
		// Update last used time
		session.mu.Lock()
		session.LastUsed = time.Now()
		session.mu.Unlock()

		// Check if session is still valid
		if session.IsHealthy() {
			logger.Trace("GetSession", "reusing existing SSH session for %s", ip)
			return session, nil
		}

		// Session is not healthy, remove it and create a new one
		logger.Trace("GetSession", "existing session for %s is unhealthy, recreating", ip)
		sm.RemoveSession(ip)
	}

	// Create new session
	return sm.CreateSession(ip, user)
}

// CreateSession creates a new SSH session for the given IP
func (sm *SessionManager) CreateSession(ip string, user string) (*SSHSession, error) {
	logger.Trace("CreateSession", "creating new SSH session for %s@%s", user, ip)

	cfg := NewSSHConfig(user, ip)
	client, err := cfg.connect()
	if err != nil {
		logger.Error("CreateSession", "failed to create SSH connection to %s: %v", ip, err)
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	now := time.Now()
	session := &SSHSession{
		Client:       client,
		Config:       cfg,
		Status:       StatusValid,
		LastUsed:     now,
		LastValidate: now,
		ConnectedAt:  now,
		FailCount:    0,
	}

	sm.mu.Lock()
	sm.sessions[ip] = session
	sm.mu.Unlock()

	logger.Log("info", "SSH session created", fmt.Sprintf("%s@%s", user, ip))
	return session, nil
}

// RemoveSession removes and closes a session
func (sm *SessionManager) RemoveSession(ip string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[ip]; exists {
		session.Close()
		delete(sm.sessions, ip)
		logger.Trace("RemoveSession", "removed SSH session for %s", ip)
	}
}

// GetSessionStatus returns the status of a session for a given IP
func (sm *SessionManager) GetSessionStatus(ip string) SessionStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if session, exists := sm.sessions[ip]; exists {
		session.mu.RLock()
		defer session.mu.RUnlock()
		return session.Status
	}

	return StatusUnknown
}

// GetAllSessionsInfo returns information about all active sessions
func (sm *SessionManager) GetAllSessionsInfo() map[string]SessionInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	info := make(map[string]SessionInfo)
	for ip, session := range sm.sessions {
		session.mu.RLock()
		info[ip] = SessionInfo{
			IP:           ip,
			Status:       session.Status,
			LastUsed:     session.LastUsed,
			LastValidate: session.LastValidate,
			ConnectedAt:  session.ConnectedAt,
			FailCount:    session.FailCount,
			Uptime:       time.Since(session.ConnectedAt),
		}
		session.mu.RUnlock()
	}

	return info
}

// SessionInfo provides information about an SSH session
type SessionInfo struct {
	IP           string        `json:"ip"`
	Status       SessionStatus `json:"status"`
	LastUsed     time.Time     `json:"last_used"`
	LastValidate time.Time     `json:"last_validate"`
	ConnectedAt  time.Time     `json:"connected_at"`
	FailCount    int           `json:"fail_count"`
	Uptime       time.Duration `json:"uptime"`
}

// validationWorker periodically validates all sessions
func (sm *SessionManager) validationWorker() {
	defer sm.wg.Done()
	ticker := time.NewTicker(sm.validationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.validateAllSessions()
		case <-sm.stopChan:
			return
		}
	}
}

// cleanupWorker periodically removes stale sessions
func (sm *SessionManager) cleanupWorker() {
	defer sm.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.cleanupStaleSessions()
		case <-sm.stopChan:
			return
		}
	}
}

// validateAllSessions validates all active sessions
func (sm *SessionManager) validateAllSessions() {
	sm.mu.RLock()
	sessions := make([]*SSHSession, 0, len(sm.sessions))
	ips := make([]string, 0, len(sm.sessions))
	for ip, session := range sm.sessions {
		sessions = append(sessions, session)
		ips = append(ips, ip)
	}
	sm.mu.RUnlock()

	for i, session := range sessions {
		ip := ips[i]
		if !session.Validate() {
			session.mu.RLock()
			failCount := session.FailCount
			session.mu.RUnlock()

			if failCount >= sm.maxFailCount {
				logger.Error("validateAllSessions", "session for %s failed %d times, removing", ip, failCount)
				sm.RemoveSession(ip)
			}
		}
	}
}

// cleanupStaleSessions removes sessions that haven't been used recently
func (sm *SessionManager) cleanupStaleSessions() {
	sm.mu.RLock()
	toRemove := make([]string, 0)
	now := time.Now()

	for ip, session := range sm.sessions {
		session.mu.RLock()
		lastUsed := session.LastUsed
		status := session.Status
		session.mu.RUnlock()

		// Remove sessions that are stale or invalid
		if now.Sub(lastUsed) > sm.staleTimeout || status == StatusInvalid || status == StatusDisconnect {
			toRemove = append(toRemove, ip)
		}
	}
	sm.mu.RUnlock()

	// Remove stale sessions
	for _, ip := range toRemove {
		logger.Trace("cleanupStaleSessions", "removing stale/invalid session for %s", ip)
		sm.RemoveSession(ip)
	}
}

// IsHealthy checks if the session is healthy without performing validation
func (s *SSHSession) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Client == nil {
		return false
	}

	if s.Status == StatusInvalid || s.Status == StatusDisconnect {
		return false
	}

	return true
}

// Validate performs a health check on the session
func (s *SSHSession) Validate() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Client == nil {
		s.Status = StatusDisconnect
		return false
	}

	// Try to create a new session and run a simple command
	session, err := s.Client.NewSession()
	if err != nil {
		s.FailCount++
		if s.FailCount >= 3 {
			s.Status = StatusDisconnect
		}
		logger.Error("Validate", "failed to create session for validation: %v", err)
		return false
	}
	defer session.Close()

	// Run a simple echo command to check if the connection is responsive
	err = session.Run("echo ok")
	if err != nil {
		s.FailCount++
		if s.FailCount >= 3 {
			s.Status = StatusInvalid
		}
		logger.Error("Validate", "validation command failed: %v", err)
		return false
	}

	// Validation successful
	s.Status = StatusValid
	s.FailCount = 0
	s.LastValidate = time.Now()
	return true
}

// Close closes the SSH session
func (s *SSHSession) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Client != nil {
		s.Client.Close()
		s.Client = nil
	}
	s.Status = StatusDisconnect
}

// RunCommand executes a command using this session
func (s *SSHSession) RunCommand(command string) (string, error) {
	s.mu.Lock()
	if s.Client == nil {
		s.mu.Unlock()
		return "", fmt.Errorf("session is closed")
	}
	client := s.Client
	s.LastUsed = time.Now()
	s.mu.Unlock()

	session, err := client.NewSession()
	if err != nil {
		s.mu.Lock()
		s.FailCount++
		s.mu.Unlock()
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command + " || [ $? = 1 ]")
	if err != nil {
		s.mu.Lock()
		s.FailCount++
		s.mu.Unlock()
		return string(output), err
	}

	// Reset fail count on success
	s.mu.Lock()
	s.FailCount = 0
	s.Status = StatusValid
	s.mu.Unlock()

	return string(output), nil
}
