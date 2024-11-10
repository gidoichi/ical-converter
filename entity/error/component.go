package error

type ComponentsError struct {
	err error
}

func NewComponentsError(err error) *ComponentsError {
	return &ComponentsError{
		err: err,
	}
}

func (e *ComponentsError) Error() string {
	return e.err.Error()
}
