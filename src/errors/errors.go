package errors

type ForbiddenError struct {
	S string
}

func (e *ForbiddenError) Error() string {
	return e.S
}
