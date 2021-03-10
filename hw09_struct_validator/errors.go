package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrOnlyStructsAllowed  = errors.New("only structs are allowed for validation")
	ErrStrLengthInvalid    = errors.New("string length invalid")
	ErrUnsupportedType     = errors.New("unsupported type")
	ErrUnknownValidator    = errors.New("unknown validator requested")
	ErrValidatorInitFailed = errors.New("failed to initialise validator")
	ErrCorruptedTag        = errors.New("given tag is corrupted")
	ErrTagParse            = errors.New("failed to parse tag")
	ErrLengthIsNotNumeric  = errors.New("failed to convert length to int")
	ErrRegexpMismatch      = errors.New("given string doesn't match pattern")
	ErrNotInChoices        = errors.New("value not present in available choices")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for _, err := range v {
		b.WriteString("field: ")
		b.WriteString(err.Field)
		b.WriteString(" ")
		b.WriteString(err.Err.Error())
		b.WriteString("\n")
	}
	return b.String()
}
