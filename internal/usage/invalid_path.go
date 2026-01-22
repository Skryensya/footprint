package usage

func InvalidPath() *Error {
	return &Error{
		Kind:    ErrInvalidPath,
		Message: "fp: path does not exist or is not accessible",
	}
}
