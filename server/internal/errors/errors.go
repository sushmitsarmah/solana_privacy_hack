package errors

import "fmt"

// ErrorResponse represents a structured error returned by the ShadowPay API.
type ErrorResponse struct {
	StatusCode   int    `json:"-"`
	Message      string `json:"message"`
	ErrorMessage string `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.ErrorMessage != "" {
		return fmt.Sprintf("shadowpay: %s (status %d) - %s", e.Message, e.StatusCode, e.ErrorMessage)
	}
	return fmt.Sprintf("shadowpay: %s (status %d)", e.Message, e.StatusCode)
}
