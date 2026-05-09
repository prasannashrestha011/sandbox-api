package sandbox_request

import "time"

// CreateRequest payload for the sanbox creation
type CreateRequest struct {
	UserID string
	// Language Environment
	Environment string
	//Docker Image ID
	ImageID string

	// -- resource limit for sandbox --
	// Memory limit
	MemoryLimit int64
	//CPU usage limit
	CPULimit int64
	//Number of processes limit
	PidsLimit int64

	// Sandbox session time expiration
	SessionTimeout time.Duration // per-user session timeout
	//Executed Command duration
	ExecTimeout time.Duration // per-command session timeout

	//Sandbox network mode
	NetWorkMode string

	CreatedAt time.Time
}
