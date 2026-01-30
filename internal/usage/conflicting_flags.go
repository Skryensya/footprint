package usage

import "fmt"

func ConflictingFlags(flag1, flag2 string) *Error {
	return &Error{
		Kind:    ErrConflictingFlags,
		Message: fmt.Sprintf("fp: flags '%s' and '%s' cannot be used together", flag1, flag2),
	}
}
