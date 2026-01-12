package usage

func FailedConfigPath() *Error {
	return &Error{
		Message:  "Getting the Config file path failed",
		ExitCode: 1,
	}
}
