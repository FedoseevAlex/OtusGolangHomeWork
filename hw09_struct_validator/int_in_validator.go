package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseIntInValidatorTag(tag string) (choices []int, err error) {
	// Assume that tag is in format in:<choices> where "choices" are comma separated allowed values
	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 || tagParts[1] == "" {
		err = ErrCorruptedTag
		return
	}

	variants := strings.Split(tagParts[1], ",")

	for _, variant := range variants {
		choice, err := strconv.Atoi(variant)
		if err != nil {
			return nil, errors.Wrap(ErrTagParse, err.Error())
		}

		choices = append(choices, choice)
	}

	return
}

type IntInValidator struct {
	BaseValidator
	Choices map[int]struct{}
	keys    []int
}

func (v *IntInValidator) Validate() {
	_, ok := v.Choices[int(v.Field.Int())]
	if ok {
		return
	}

	v.Errs = append(
		v.Errs,
		ValidationError{
			Field: v.FieldName,
			Err: errors.Wrapf(
				ErrNotInChoices,
				"value: %s, is not in %v",
				v.Field.String(),
				v.keys,
			),
		},
	)
}

func NewIntInValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	numbers, err := parseIntInValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	choices := make(map[int]struct{}, len(numbers))
	for _, number := range numbers {
		choices[number] = struct{}{}
	}

	return &IntInValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Choices:       choices,
		keys:          numbers,
	}, nil
}

type IntInSliceValidator IntInValidator

func (v *IntInSliceValidator) Validate() {
	for i := 0; i < v.Field.Len(); i++ {
		num := int(v.Field.Index(i).Int())
		_, ok := v.Choices[num]
		if ok {
			return
		}

		v.Errs = append(
			v.Errs,
			ValidationError{
				Field: v.FieldName,
				Err: errors.Wrapf(
					ErrNotInChoices,
					"value: %d, is not in %v",
					v.Field.Int(),
					v.keys,
				),
			},
		)
	}
}

func NewIntInSliceValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	numbers, err := parseIntInValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	choices := make(map[int]struct{}, len(numbers))
	for _, number := range numbers {
		choices[number] = struct{}{}
	}

	return &IntInSliceValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Choices:       choices,
		keys:          numbers,
	}, nil
}
