package dto

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorBody  `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}
type ErrorBody struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}
