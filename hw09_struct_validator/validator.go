package hw09structvalidator

import "reflect"

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

	for i := 0; i < vVal.NumField(); i++ {
		fieldVal := vVal.Field(i)
		if !fieldVal.CanInterface() {
			// skip unexported field
			continue
		}

		fieldInfo := vType.Field(i)
		if _, ok := fieldInfo.Tag.Lookup(validateTagName); !ok {
			// skip if validation tag not found
			continue
		}

		vs, err := prepareValidators(fieldInfo)
		if err != nil {
			return err
		}

		for _, v := range vs {
			errs := v.IsValid(fieldVal.Interface())
			if len(errs) != 0 {
				vErrs = append(vErrs, errs...)
			}
		}
	}

	if len(vErrs) != 0 {
		return vErrs
	}
	return nil
}

// This function parses tag and calls initializers for validators
func prepareValidators(field reflect.StructField) ([]Validator, error) {
	panic("implement me")
	return []Validator{}, nil
}
