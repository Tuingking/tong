package gotong

import "fmt"

type BuilderError struct {
	Where string
	Why   string
}

func (e *BuilderError) Error() error {
	return fmt.Errorf("there is an error at %s due to %s", e.Where, e.Why)
}

func NewBuilderError(where string, why error) BuilderError {
	return BuilderError{
		Where: where,
		Why:   why.Error(),
	}
}
