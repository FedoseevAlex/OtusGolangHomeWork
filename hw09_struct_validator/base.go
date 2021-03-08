package hw09_struct_validator //nolint:golint,stylecheck

import "reflect"

type Validator interface {
	Validate()
}

type BaseValidator struct {
	Field reflect.Value
	Errs  ValidationErrors
}

// Type for validator initializer function.
// Requires reflect.Value to validate it
type ValidatorInitFunc func(reflect.Value, reflect.StructField) (*Validator, error)
