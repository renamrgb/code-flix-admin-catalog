// Package validation provides error types and utilities for domain validation.
package validation

import "errors"

type ValidationErrors struct {
	Errs []error
}

func (v ValidationErrors) Error() string {
	return errors.Join(v.Errs...).Error()
}
