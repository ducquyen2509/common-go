package errors

import "fmt"

type ParseError struct {
	Param      string
	ExpectType string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse error param %s, expect type %s", e.Param, e.ExpectType)
}
