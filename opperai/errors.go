package opperai

import "errors"

var (
	ErrRateLimit       = errors.New("rate limit error: please retry in a few seconds")
	ErrFunctionRunFail = errors.New("failed to run function")
)
