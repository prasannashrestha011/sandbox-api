package enums

// SandboxState represents the lifecycle state of a sandbox.
type SandboxState string

const (
	StateActive   SandboxState = "active"
	StateInActive SandboxState = "inactive"
)

// User roles
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type UserType string

const (
	UserTypeStudent    UserType = "student"
	UserTypeInstructor UserType = "instructor"
)

type LabProgress string

const (
	Enrolled           LabProgress = "enrolled"
	ProgressInProgress LabProgress = "in_progress"
	ProgressCompleted  LabProgress = "completed"
)

type SubmissionStatus string

const (
	SubmissionAccepted SubmissionStatus = "accepted"
	SubmissionRejected SubmissionStatus = "rejected"
)
