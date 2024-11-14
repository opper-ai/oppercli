package opperai

import "errors"

var (
	ErrRateLimit       = errors.New("rate limit error: please retry in a few seconds")
	ErrFunctionRunFail = errors.New("failed to run function")
	ErrUnauthorized    = errors.New("unauthorized: invalid API key")
	ErrNotFound        = errors.New("not found")
)

// IsErrorType checks if an error is of a specific type
func IsErrorType(err, target error) bool {
	return errors.Is(err, target)
}
