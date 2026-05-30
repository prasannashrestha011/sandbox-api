package dto

type FieldViolation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Violations []FieldViolation
}

func (v ValidationErrors) Error() string {
	return "validation failed"
}
