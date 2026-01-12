package usage

import "fmt"

func InvalidFlag(flag string) *Error {
	return &Error{
		Message:  fmt.Sprintf("fp: invalid flag '%s'", flag),
		ExitCode: 2,
	}
}
