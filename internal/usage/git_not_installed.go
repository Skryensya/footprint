package usage

func GitNotInstalled() *Error {
	return &Error{
		Kind:    ErrGitNotInstalled,
		Message: "fp: git command not found in PATH",
	}
}
