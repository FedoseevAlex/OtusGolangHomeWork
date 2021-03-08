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

func (v *StrLenValidator) Validate() {
	if v.Field.Len() != v.Len {
		return
	}

	v.Errs = append(
		v.Errs,
		ValidationError{
			Field: v.Field.Type().Name(),
			Err:   errors.Wrapf(ErrStrLengthInvalid, "need %d, got %d", v.Len, v.Field.Len()),
		},
	)
}

func NewStrLenValidator(fieldValue reflect.Value, fieldInfo reflect.StructField) (*StrLenValidator, error) {
	// Assume that tag is in format len:<N> where N is desired length
	tag := fieldInfo.Tag.Get(validateTagName)
	fieldName := fieldInfo.Name

	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 {
		return nil, errors.Wrapf(ErrCorruptedTag, "Field: %s Tag: %s", fieldName, tag)
	}

	length, err := strconv.Atoi(tagParts[1])
	if err != nil {
		return nil, errors.Wrapf(ErrLengthIsNotNumeric, "Field: %s Tag: %s", fieldName, tag)
	}

	return &StrLenValidator{
		BaseValidator: BaseValidator{Field: fieldValue},
		Len:           length,
	}, nil
}

type StrLenSliceValidator StrLenValidator

func (v StrLenSliceValidator) IsValid() {
	for i := 0; i < v.Field.Len(); i++ {
		elem := v.Field.Index(i)
		if elem.Len() != v.Len {
			return
		}

		v.Errs = append(
			v.Errs,
			ValidationError{
				Field: v.Field.Type().Name(),
				Err:   errors.Wrapf(ErrStrLengthInvalid, "need %d, got %d", v.Len, v.Field.Len()),
			},
		)
	}
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
