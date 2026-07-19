package pictomancer

import "fmt"

// APIError is returned for any non-2xx response from the API.
type APIError struct {
	Status int
	Detail string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("pictomancer: HTTP %d: %s", e.Status, e.Detail)
}

func newAPIError(status int, detail string) *APIError {
	return &APIError{Status: status, Detail: detail}
}
