package usage

func NotInGitRepo() *Error {
	return &Error{
		Message:  "the path is not inside a valid git repository",
		ExitCode: 1,
	}
}
