package hw09_struct_validator //nolint:golint,stylecheck

type Validator interface {
	IsValid(interface{}) ValidationErrors
}

type BaseValidator struct {
	Field string
}
