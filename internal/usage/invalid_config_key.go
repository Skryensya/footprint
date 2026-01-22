package usage

import "fmt"

func InvalidConfigKey(key string) *Error {
	return &Error{
		Kind:    ErrInvalidConfigKey,
		Message: fmt.Sprintf("fp: unknown config key '%s'", key),
	}
}
