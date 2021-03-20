package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// TODO: write test here.
func parseIntMinValidatorTag(tag string) (min int, err error) {
	// Assume that tag is in format min:<N> where "N" is numeric string
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

type IntMinValidator struct {
	BaseValidator
	Min int
}

func (v *IntMinValidator) Validate() {
	value := v.Field.Int()
	if int(value) >= v.Min {
		// It's okay if value greater or equal to minimum.
		return
	}

	v.Errs = append(
		v.Errs,
		ValidationError{
			Field: v.FieldName,
			Err: errors.Wrapf(
				ErrBelowMinimum,
				"value: %d < minimum: %d",
				value,
				v.Min,
			),
		},
	)
}

func NewIntMinValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	minimum, err := parseIntMinValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &IntMinValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Min:           minimum,
	}, nil
}

type IntMinSliceValidator IntMinValidator

func (v *IntMinSliceValidator) Validate() {
	for i := 0; i < v.Field.Len(); i++ {
		value := int(v.Field.Index(i).Int())

		if int(value) >= v.Min {
			// It's okay if value greater or equal to minimum.
			return
		}

		v.Errs = append(
			v.Errs,
			ValidationError{
				Field: v.FieldName,
				Err: errors.Wrapf(
					ErrBelowMinimum,
					"value: %d < minimum: %d",
					value,
					v.Min,
				),
			},
		)
	}
}

func NewIntMinSliceValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	minimum, err := parseIntMinValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &IntMinSliceValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Min:           minimum,
	}, nil
}
