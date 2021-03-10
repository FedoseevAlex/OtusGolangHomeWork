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

// TODO: write test here.
func parseLenValidatorTag(tag string) (length int, err error) {
	// Assume that tag is in format len:<N> where N is desired length
	tagParts := strings.Split(tag, ":")
	if len(tagParts) != 2 {
		err = ErrCorruptedTag
		return
	}

	length, err = strconv.Atoi(tagParts[1])
	if err != nil {
		return
	}

	return
}

func NewStrLenValidator(fieldValue reflect.Value, fieldInfo reflect.StructField) (Validator, error) {
	tag := fieldInfo.Tag.Get(validateTagName)
	fieldName := fieldInfo.Name

	length, err := parseLenValidatorTag(tag)
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

func NewStrLenSliceValidator(fieldValue reflect.Value, fieldInfo reflect.StructField) (Validator, error) {
	tag := fieldInfo.Tag.Get(validateTagName)
	fieldName := fieldInfo.Name

	length, err := parseLenValidatorTag(tag)
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
