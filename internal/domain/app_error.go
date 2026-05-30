package domain

type AppError struct {
	HTTPCode int
	Code     ErrorCode
	Message  string
	// Internal error (never exposed)
	Err error `json:"-"`

	// Optional structured details (safe to expose)
	ErrorDetails any `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
