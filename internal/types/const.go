package types

// SandboxState represents the lifecycle state of a sandbox.
type SandboxState string

const (
	StateActive   SandboxState = "active"
	StateInActive SandboxState = "inactive"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)
