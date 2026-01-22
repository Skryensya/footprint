package usage

import "fmt"

func InvalidFlag(flag string) *Error {
	return &Error{
		Kind:    ErrInvalidFlag,
		Message: fmt.Sprintf("fp: invalid flag '%s'", flag),
	}
}
