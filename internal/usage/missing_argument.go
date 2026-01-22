package usage

import "fmt"

func MissingArgument(arg string) *Error {
	return &Error{
		Kind:    ErrMissingArgument,
		Message: fmt.Sprintf("fp: missing required argument '%s'", arg),
	}
}
