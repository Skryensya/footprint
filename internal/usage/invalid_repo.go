package usage

func InvalidRepo() *Error {
	return &Error{
		Kind:    ErrInvalidRepo,
		Message: "fp: could not identify repository (no remote URL or valid path)",
	}
}
