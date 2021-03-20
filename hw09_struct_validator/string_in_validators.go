package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func parseStrInValidatorTag(tag string) (choices []string, err error) {
	// Assume that tag is in format in:<choices> where "choices" are comma separated allowed values
	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 || tagParts[1] == "" {
		err = ErrCorruptedTag
		return
	}

	choices = strings.Split(tagParts[1], ",")
	return
}

type StrInValidator struct {
	BaseValidator
	Choices map[string]struct{}
	keys    []string
}

func (v *StrInValidator) Validate() {
	_, ok := v.Choices[v.Field.String()]
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

func NewStrInValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	strs, err := parseStrInValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	choices := make(map[string]struct{}, len(strs))
	for _, str := range strs {
		choices[str] = struct{}{}
	}

	return &StrInValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Choices:       choices,
		keys:          strs,
	}, nil
}

type StrInSliceValidator StrInValidator

func (v *StrInSliceValidator) Validate() {
	for i := 0; i < v.Field.Len(); i++ {
		str := v.Field.Index(i).String()

		_, ok := v.Choices[str]
		if ok {
			continue
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
}

func NewStrInSliceValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	strs, err := parseStrInValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	choices := make(map[string]struct{}, len(strs))
	for _, str := range strs {
		choices[str] = struct{}{}
	}

	return &StrInSliceValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Choices:       choices,
		keys:          strs,
	}, nil
}
