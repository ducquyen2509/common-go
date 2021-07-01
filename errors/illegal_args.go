package errors

import "fmt"

type IllegalArgumentsError struct {
	Args []string
}

func (e IllegalArgumentsError) Error() string {
	return fmt.Sprintf("illegal arguments %v", e.Args)
}
