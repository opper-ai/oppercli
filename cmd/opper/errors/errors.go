package errors

import "fmt"

func NotFound(resource, name string) error {
	return fmt.Errorf("%s not found: %s", resource, name)
}

func InvalidInput(msg string) error {
	return fmt.Errorf("invalid input: %s", msg)
}
