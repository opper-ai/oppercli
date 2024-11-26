package builders

import (
	"fmt"
	"net"
	"strings"
)

// FormatError formats errors in a user-friendly way
func FormatError(err error) error {
	if err == nil {
		return nil
	}

	// Check for connection errors
	if netErr, ok := err.(*net.OpError); ok {
		if netErr.Op == "dial" {
			return fmt.Errorf("connection failed: could not connect to server. Is the server running?")
		}
		return fmt.Errorf("network error: %v", netErr)
	}

	// Remove any trailing newlines from error messages
	msg := strings.TrimSpace(err.Error())

	// If it's already a formatted error, return as is
	if strings.Contains(msg, "error:") {
		return fmt.Errorf(msg)
	}

	return fmt.Errorf("error: %v", msg)
}
