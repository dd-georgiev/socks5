package command_response

import "fmt"

type InvalidAtypError struct {
	Atyp uint16
}

func (e *InvalidAtypError) Error() string {
	return fmt.Sprintf("Invalid Atyp: %d", e.Atyp)
}

type InvalidStatusError struct {
	Status uint16
}

func (e *InvalidStatusError) Error() string {
	return fmt.Sprintf("Invalid Status: %d", e.Status)
}
