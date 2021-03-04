package hw09_struct_validator //nolint:golint,stylecheck

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type StrLenValidator struct {
	BaseValidator
	Len int
}

func (v StrLenValidator) IsValid(s interface{}) ValidationErrors {
	vVal := reflect.ValueOf(s)
	if vVal.Len() != v.Len {
		return ValidationErrors{
			ValidationError{
				Field: v.Field,
				Err:   errors.Wrapf(ErrStrLengthInvalid, "need %d, got %d", v.Len, vVal.Len()),
			}}
	}
	return nil
}

func NewStrLenValidator(fieldName, tag string) (StrLenValidator, error) {
	// Assume that tag is in format len:<N> where N is desired length
	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 {
		return StrLenValidator{}, errors.Wrapf(ErrCorruptedTag, "Field: %s Tag: %s", fieldName, tag)
	}

	length, err := strconv.Atoi(tagParts[1])
	if err != nil {
		return StrLenValidator{}, errors.Wrapf(ErrLengthIsNotNumeric, "Field: %s Tag: %s", fieldName, tag)
	}

	return StrLenValidator{
		BaseValidator: BaseValidator{Field: fieldName},
		Len:           length,
	}, nil
}

type StrLenSliceValidator StrLenValidator

func (v StrLenSliceValidator) IsValid(s interface{}) ValidationErrors {
	sVal := reflect.ValueOf(s)

	errs := make(ValidationErrors, 0, sVal.Len())
	for i := 0; i < sVal.Len(); i++ {
		err := v.IsValid(sVal.Index(i))
		if err != nil {
			errs = append(errs, err...)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func NewStrLenSliceValidator(fieldName, tag string) (StrLenSliceValidator, error) {
	sv, err := NewStrLenValidator(fieldName, tag)
	return StrLenSliceValidator(sv), err
}

type StrRegexpValidator struct {
	BaseValidator
	pattern *regexp.Regexp
}

func NewStrRegexpValidator(fieldName, tag string) (StrRegexpValidator, error) {
	// Assume that tag is like "regexp:<pattern>"
	tagParts := strings.SplitN(tag, ":", 1)
	if len(tagParts) != 2 {
		return StrRegexpValidator{}, errors.Wrapf(ErrCorruptedTag, "Field: %s Tag: %s", fieldName, tag)
	}

	pattern, err := regexp.Compile(tagParts[1])
	if err != nil {
		return StrRegexpValidator{}, errors.WithMessagef(err, "Field: %s Tag: %s", fieldName, tag)
	}

	return StrRegexpValidator{BaseValidator: BaseValidator{Field: fieldName}, pattern: pattern}, nil
}
