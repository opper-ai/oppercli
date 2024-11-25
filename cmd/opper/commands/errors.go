package commands

import "fmt"

type CommandError struct {
	Command string
	Action  string
	Err     error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("%s command %s failed: %v", e.Command, e.Action, e.Err)
}

func WrapError(command, action string, err error) error {
	return &CommandError{
		Command: command,
		Action:  action,
		Err:     err,
	}
}
