package services_validators

import (
	"main/internal/services/models"
)

const (
	MaxMemoryLimit = 512 * 1024 * 1024 // 512 MB
	MaxCPULimit    = 2                 // 2 CPUs
	MaxPidsLimit   = 100               // 100 Processes
)

// ValidateAndCapSandboxLimits enforces sensible resource limits, preventing
// malicious or excessive resource allocation.
func ValidateAndCapSandboxLimits(sandbox *models.SandboxTemplate) {
	if sandbox == nil {
		return
	}

	if sandbox.MemoryLimit <= 0 || sandbox.MemoryLimit > MaxMemoryLimit {
		sandbox.MemoryLimit = MaxMemoryLimit
	}
	if sandbox.CPULimit <= 0 || sandbox.CPULimit > MaxCPULimit {
		sandbox.CPULimit = MaxCPULimit
	}
	if sandbox.PidsLimit <= 0 || sandbox.PidsLimit > MaxPidsLimit {
		sandbox.PidsLimit = MaxPidsLimit
	}
}
