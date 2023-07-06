package app_errors

import "fmt"

type DBError struct {
	Ty  string
	Err error
}

type InternalError DBError

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error (%s): %v", e.Ty, e.Err.Error())
}

type ValueError DBError

func (e *ValueError) Error() string {
	return fmt.Sprintf("value error (%s): %v", e.Ty, e.Err.Error())
}

type NotFoundError DBError

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found error (%s): %v", e.Ty, e.Err.Error())
}
