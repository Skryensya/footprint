package usage

func InvalidRepo() *Error {
	return &Error{
		Message:  "Invalid repository",
		ExitCode: 1,
	}
}
