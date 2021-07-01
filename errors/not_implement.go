package errors

import "fmt"

type NotImplementedError struct {
	Func string
}

func (e NotImplementedError) Error() string {
	return fmt.Sprintf("not implement func %s", e.Func)
}
