package usage

func MissingRemote() *Error {
	return &Error{
		Kind:    ErrMissingRemote,
		Message: "fp: current repository does not have a remote",
	}
}
