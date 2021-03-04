package hw09structvalidator

import "reflect"

func Validate(v interface{}) error {
	// Place your code here.
	return nil
}

func prepareValidator(field reflect.StructField) ([]Validator, error) {
	panic("implement me")
	return []Validator{}, nil
}
