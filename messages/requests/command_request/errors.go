package command_request

import "fmt"

type InvalidCommandError struct {
	CommandType uint16
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("Invalid Command: %d", e.CommandType)
}
