package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Parses tag and tries to compile pattern.
func parseRegexpValidatorTag(tag string) (pattern *regexp.Regexp, err error) {
	tagParts := strings.SplitN(tag, ":", 2)
	if len(tagParts) < 2 || tagParts[1] == "" {
		err = ErrCorruptedTag
		return
	}

	pattern, err = regexp.Compile(tagParts[1])
	if err != nil {
		err = errors.Wrap(ErrTagParse, err.Error())
		pattern = nil
		return
	}

	return
}

type StrRegexpValidator struct {
	BaseValidator
	Pattern *regexp.Regexp
}

func (v *StrRegexpValidator) Validate() {
	s := v.Field.String()

	ok := v.Pattern.MatchString(s)
	if !ok {
		v.Errs = append(v.Errs,
			ValidationError{
				Field: v.Field.Type().Name(),
				Err:   errors.WithMessagef(ErrRegexpMismatch, "string: %s, pattern: %s", s, v.Pattern.String()),
			},
		)
	}
}

func NewStrRegexpValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	pattern, err := parseRegexpValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &StrRegexpValidator{
		BaseValidator: BaseValidator{Field: fieldValue}, Pattern: pattern,
	}, nil
}

type StrRegexpSliceValidator StrRegexpValidator

func NewStrRegexpSliceValidator(fieldValue reflect.Value, fieldName, tag string) (Validator, error) {
	pattern, err := parseRegexpValidatorTag(tag)
	if err != nil {
		return nil, errors.Wrapf(
			ErrTagParse,
			"field: %s, tag: %s, err: %s",
			fieldName,
			tag,
			err.Error(),
		)
	}

	return &StrRegexpSliceValidator{
		BaseValidator: BaseValidator{Field: fieldValue}, Pattern: pattern,
	}, nil
}

func (v *StrRegexpSliceValidator) Validate() {
	for i := 0; i < v.Field.Len(); i++ {
		s := v.Field.Index(i)
		ok := v.Pattern.MatchString(s.String())
		if !ok {
			v.Errs = append(v.Errs,
				ValidationError{
					Field: v.Field.Type().Name(),
					Err:   errors.WithMessagef(ErrRegexpMismatch, "string: %s, pattern: %s", s, v.Pattern.String()),
				},
			)
		}
	}
}
