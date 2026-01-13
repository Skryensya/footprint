package usage

func MissingRemote() *Error {
	return &Error{
		Message:  "Current repository does not have a remote",
		ExitCode: 2,
	}
}
