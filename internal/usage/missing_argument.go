package usage

import "fmt"

func MissingArgument(arg string) *Error {
	return &Error{
		Message:  fmt.Sprintf("fp: missing required argument '%s'", arg),
		ExitCode: 2,
	}
}
