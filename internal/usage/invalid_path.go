package usage

func InvalidPath() *Error {
	return &Error{
		Message:  "Invalid path",
		ExitCode: 1,
	}
}
