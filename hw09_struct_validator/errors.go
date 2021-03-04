package hw09_struct_validator //nolint:golint,stylecheck

import (
	"github.com/pkg/errors"
)

var (
	ErrOnlyStructsAllowed = errors.New("only structs are allowed for validation")
	ErrStrLengthInvalid   = errors.New("string length invalid")
	ErrUnsupportedType    = errors.New("unsupported type")

	ErrCorruptedTag       = errors.New("given tag is corrupted")
	ErrLengthIsNotNumeric = errors.New("failed to convert length to int")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}
