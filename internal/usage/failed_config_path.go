package usage

func FailedConfigPath() *Error {
	return &Error{
		Kind:    ErrFailedConfigPath,
		Message: "fp: could not determine config file location",
	}
}
