package hw09_struct_validator //nolint:golint,stylecheck,revive

import "reflect"

type Validator interface {
	Validate()
	Errors() ValidationErrors
}

type BaseValidator struct {
	Field     reflect.Value
	FieldName string
	Errs      ValidationErrors
}

func (v BaseValidator) Errors() ValidationErrors {
	return v.Errs
}

// Type for validator initializer function.
// Requires reflect.Value and reflect.StructField to build a validator.
type ValidatorInitFunc func(fieldValue reflect.Value, name, tag string) (Validator, error)

type ValidationSelectorKey struct {
	Kind           reflect.Kind
	ElemKind       reflect.Kind
	ValidationType string
}

type ValidationSelector map[ValidationSelectorKey]ValidatorInitFunc
