package commands

// UsageError represents an error that should show command usage
type UsageError struct {
	err error
}

func NewUsageError(err error) *UsageError {
	return &UsageError{err: err}
}

func (e *UsageError) Error() string {
	return e.err.Error()
}

// IsUsageError checks if an error should show usage
func IsUsageError(err error) bool {
	_, ok := err.(*UsageError)
	return ok
}
