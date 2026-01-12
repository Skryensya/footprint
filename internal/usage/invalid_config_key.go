package usage

import "fmt"

func InvalidConfigKey(key string) *Error {
	return &Error{
		Message:  fmt.Sprintf("fp: there is not config with key '%s'", key),
		ExitCode: 1,
	}
}
