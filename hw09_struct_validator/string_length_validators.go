package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type StrLenValidator struct {
	BaseValidator
	Len int
}

func (v *StrLenValidator) Validate() {
	if v.Field.Len() == v.Len {
		return
	}

	v.Errs = append(
		v.Errs,
		ValidationError{
			Field: v.FieldName,
			Err: errors.Wrapf(
				ErrStrLengthInvalid,
				"value: %s, need %d, got %d",
				v.Field.String(),
				v.Len,
				v.Field.Len(),
			),
		},
	)
}

func parseStrLenValidatorTag(tag string) (length int, err error) {
	// Assume that tag is in format len:<N> where N is desired length
	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 || tagParts[1] == "" {
		err = ErrCorruptedTag
		return
	}

	length, err = strconv.Atoi(tagParts[1])
	if err != nil {
		err = errors.Wrap(ErrTagParse, err.Error())
		return
	}

	if length < 0 {
		err = errors.Wrapf(ErrTagParse, "negative length specified")
		length = 0
		return
	}

	return
}

func NewStrLenValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	length, err := parseStrLenValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &StrLenValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Len:           length,
	}, nil
}

type StrLenSliceValidator StrLenValidator

func (v *StrLenSliceValidator) Validate() {
	for i := 0; i < v.Field.Len(); i++ {
		elem := v.Field.Index(i)
		if elem.Len() == v.Len {
			continue
		}

		v.Errs = append(
			v.Errs,
			ValidationError{
				Field: v.FieldName,
				Err: errors.Wrapf(
					ErrStrLengthInvalid,
					"value: %s, need %d, got %d",
					elem.String(),
					v.Len,
					elem.Len(),
				),
			},
		)
	}
}

func NewStrLenSliceValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	length, err := parseStrLenValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &StrLenSliceValidator{
		BaseValidator: BaseValidator{Field: fieldValue, FieldName: fieldName},
		Len:           length,
	}, nil
}
