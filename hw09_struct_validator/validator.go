package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const (
	validateTagName = "validate"
)

func Validate(v interface{}) error {
	vType := reflect.TypeOf(v)
	if vType.Kind() != reflect.Struct {
		return ErrOnlyStructsAllowed
	}
	vVal := reflect.ValueOf(v)

	vErrs := make(ValidationErrors, 0, vVal.NumField())
	validators := make([]Validator, 0, vVal.NumField())

	for i := 0; i < vVal.NumField(); i++ {
		fieldVal := vVal.Field(i)
		fieldName := vType.Field(i).Name
		fieldTag := vType.Field(i).Tag

		_ = fieldName
		_ = fieldTag

		if !fieldVal.CanInterface() {
			// skip unexported field
			continue
		}

		fieldInfo := vType.Field(i)
		if _, ok := fieldInfo.Tag.Lookup(validateTagName); !ok {
			// skip if validation tag not found
			continue
		}

		vs, err := prepareValidators(fieldVal, fieldInfo)
		if err != nil {
			return err
		}

		validators = append(validators, vs...)
	}

	for _, v := range validators {
		v.Validate()
		vErrs = append(vErrs, v.Errors()...)
	}

	if len(vErrs) != 0 {
		return vErrs
	}
	return nil
}

var validationSelector = ValidationSelector{
	{
		Kind:           reflect.String,
		ValidationType: "len",
	}: NewStrLenValidator,
	{
		Kind:           reflect.Slice,
		ElemKind:       reflect.String,
		ValidationType: "len",
	}: NewStrLenSliceValidator,
	{
		Kind:           reflect.String,
		ValidationType: "regexp",
	}: NewStrRegexpValidator,
	{
		Kind:           reflect.Slice,
		ElemKind:       reflect.String,
		ValidationType: "regexp",
	}: NewStrRegexpSliceValidator,
}

// This function parses tag and calls initializers for validators
func prepareValidators(field reflect.Value, fieldInfo reflect.StructField) (vs []Validator, err error) {
	var (
		kind     reflect.Kind
		elemKind reflect.Kind
	)

	kind = fieldInfo.Type.Kind()
	if kind == reflect.Slice {
		elemKind = fieldInfo.Type.Elem().Kind()
	}

	// Validator descriptions
	vds := fieldInfo.Tag.Get(validateTagName)
	vdsParts := strings.Split(vds, "|")

	for _, vd := range vdsParts {
		vType := strings.SplitN(vd, ":", 2)
		key := ValidationSelectorKey{
			Kind:           kind,
			ElemKind:       elemKind,
			ValidationType: vType[0],
		}

		initFunc, ok := validationSelector[key]
		if !ok {
			err = errors.Wrapf(ErrUnknownValidator, "validator not found: %s %s", vType, vd)
			return
		}

		var validator Validator

		validator, err = initFunc(field, fieldInfo)
		if err != nil {
			return
		}

		vs = append(vs, validator)
	}

	return
}
