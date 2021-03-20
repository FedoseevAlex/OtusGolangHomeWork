package hw09_struct_validator //nolint:golint,stylecheck,revive,dupl

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseIntMaxValidatorTag(tag string) (min int, err error) {
	// Assume that tag is in format max:<N> where "N" is numeric string
	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 || tagParts[1] == "" {
		err = ErrCorruptedTag
		return
	}

	min, err = strconv.Atoi(tagParts[1])
	if err != nil {
		return 0, errors.Wrap(ErrTagParse, err.Error())
	}

	return
}

type IntMaxValidator struct {
	BaseValidator
	Max int
}

func (v *IntMaxValidator) Validate() {
	value := v.Field.Int()
	if int(value) <= v.Max {
		// It's okay if value less or equal to maximum.
		return
	}

	v.Errs = append(
		v.Errs,
		ValidationError{
			Field: v.FieldName,
			Err: errors.Wrapf(
				ErrAboveMaximum,
				"value: %d > maximum: %d",
				value,
				v.Max,
			),
		},
	)
}

func NewIntMaxValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	maximum, err := parseIntMaxValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &IntMaxValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Max:           maximum,
	}, nil
}

type IntMaxSliceValidator IntMaxValidator

func (v *IntMaxSliceValidator) Validate() {
	for i := 0; i < v.Field.Len(); i++ {
		value := int(v.Field.Index(i).Int())

		if value <= v.Max {
			// It's okay if value less or equal to maximum.
			return
		}

		v.Errs = append(
			v.Errs,
			ValidationError{
				Field: v.FieldName,
				Err: errors.Wrapf(
					ErrAboveMaximum,
					"value: %d > maximum: %d",
					value,
					v.Max,
				),
			},
		)
	}
}

func NewIntMaxSliceValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	maximum, err := parseIntMaxValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &IntMaxSliceValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Max:           maximum,
	}, nil
}
