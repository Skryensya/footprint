package usage

func GitNotInstalled() *Error {
	return &Error{
		Message:  "Git is not installed in the system",
		ExitCode: 1,
	}
}
