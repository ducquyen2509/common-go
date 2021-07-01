package errors

import "fmt"

type MissingArgumentsError struct {
	Args []string
}

func (e MissingArgumentsError) Error() string {
	return fmt.Sprintf("missing arguments %v", e.Args)
}
