package usage

func NotInGitRepo() *Error {
	return &Error{
		Kind:    ErrNotInGitRepo,
		Message: "fp: path is not inside a valid git repository",
	}
}
