package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type User struct {
	ID     string `json:"id" validate:"len:36"`
	Name   string
	Age    int      `validate:"min:18|max:50"`
	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$|len:19"`
	Role   UserRole `validate:"in:admin,stuff"`
	Phones []string `validate:"len:11"`
	meta   json.RawMessage
}

func TestUserValidators(t *testing.T) {
	tests := []struct {
		name             string
		in               interface{}
		expectedErrCount int
		expectedErr      error
	}{
		{
			name: "user struct valid",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    33,
				Email:  "johnny_doe@yolo.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
		},
		{
			name: "user struct invalid id",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2-dead-face-cafe",
				Name:   "John Doe",
				Age:    33,
				Email:  "johnny_doe@yolo.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 1,
			expectedErr:      ErrStrLengthInvalid,
		},
		{
			name: "user struct invalid age: too young",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    10,
				Email:  "johnny_doe@yolo.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 1,
			expectedErr:      ErrBelowMinimum,
		},
		{
			name: "user struct invalid age: too old",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    100,
				Email:  "johnny_doe@yolo.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 1,
			expectedErr:      ErrAboveMaximum,
		},
		{
			name: "user struct invalid email: regexp mismatch",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    33,
				Email:  "johnny_doe#yolo.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 1,
			expectedErr:      ErrRegexpMismatch,
		},
		{
			name: "user struct invalid email: length mismatch",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    33,
				Email:  "john_doe@yolo.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 1,
			expectedErr:      ErrStrLengthInvalid,
		},
		{
			name: "user struct invalid role",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    33,
				Email:  "johnny_doe@yolo.com",
				Role:   "CEO",
				Phones: []string{"00000000000"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 1,
			expectedErr:      ErrNotInChoices,
		},
		{
			name: "user struct invalid phone",
			in: User{
				ID:     "1ebe928b-5044-43bb-a8b6-783a7dec48c2",
				Name:   "John Doe",
				Age:    33,
				Email:  "johnny_doe@yolo.com",
				Role:   "admin",
				Phones: []string{"12345678901", "12345678901234", "123456"},
				meta:   json.RawMessage{},
			},
			expectedErrCount: 2,
			expectedErr:      ErrStrLengthInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, errs)
				return
			}

			var validationErrors ValidationErrors
			require.ErrorAs(t, errs, &validationErrors)
			require.Equal(t, tt.expectedErrCount, len(validationErrors))

			for _, err := range validationErrors {
				require.ErrorIs(t, err.Err, tt.expectedErr)
			}
		})
	}
}

// Here we check basic functionality of every validator.
// We are interested only in validation errors so leave aside
// errors that may occur at validator preparation.
func TestValidatorsBasic(t *testing.T) {
	tests := []struct {
		name             string
		in               interface{}
		expectedErrCount int
		expectedErr      error
	}{
		{
			name: "valid length case",
			in: struct {
				Value string `validate:"len:5"`
			}{Value: "five!"},
		}, {
			name: "invalid length case",
			in: struct {
				Value string `validate:"len:5"`
			}{Value: "definetly not five"},
			expectedErrCount: 1,
			expectedErr:      ErrStrLengthInvalid,
		},
		{
			name: "valid length case in slice",
			in: struct {
				Value []string `validate:"len:2"`
			}{[]string{"aa", "bb"}},
		},
		{
			name: "invalid length case in slice",
			in: struct {
				Value []string `validate:"len:2"`
			}{[]string{"a", "bb", "ccc"}},
			expectedErrCount: 2,
			expectedErr:      ErrStrLengthInvalid,
		},
		{

			name: "valid regexp check",
			in: struct {
				Value string `validate:"regexp:^a+$"`
			}{"aa"},
		},
		{
			name: "invalid regexp check",
			in: struct {
				Value string `validate:"regexp:^a+$"`
			}{"adfasdkfj;aslkd12341234"},
			expectedErrCount: 1,
			expectedErr:      ErrRegexpMismatch,
		},
		{
			name: "valid regexp check in slice",
			in: struct {
				Value []string `validate:"regexp:^\\d+$"`
			}{[]string{"12312", "34234234234234"}},
		},
		{
			name: "invalid regexp check in slice",
			in: struct {
				Value []string `validate:"regexp:^\\d+$"`
			}{[]string{"123", "0xDEADFACECAFE", ""}},
			expectedErrCount: 2,
			expectedErr:      ErrRegexpMismatch,
		},

		{
			name: "valid 'in' check",
			in: struct {
				Value string `validate:"in:a,b,c,d"`
			}{"a"},
		},
		{
			name: "failed 'in' check",
			in: struct {
				Value string `validate:"in:a,b,c,d"`
			}{"e"},
			expectedErrCount: 1,
			expectedErr:      ErrNotInChoices,
		},
		{
			name: "valid 'in' check in slice",
			in: struct {
				Value []string `validate:"in:a,b,c,d"`
			}{[]string{"a", "b", "c"}},
		},
		{
			name: "failed 'in' check in slice",
			in: struct {
				Value []string `validate:"in:a,b,c,d"`
			}{[]string{"e", "g", "c"}},
			expectedErrCount: 2,
			expectedErr:      ErrNotInChoices,
		},
		{
			name: "valid 'min' check",
			in: struct {
				Value int `validate:"min:10"`
			}{Value: 11},
		},
		{
			name: "failed 'min' check",
			in: struct {
				Value int `validate:"min:10"`
			}{Value: -1},
			expectedErrCount: 1,
			expectedErr:      ErrBelowMinimum,
		},
		{
			name: "valid 'min' check slice",
			in: struct {
				Value []int `validate:"min:10"`
			}{Value: []int{10, 11, 12}},
		},
		{
			name: "failed 'min' check slice",
			in: struct {
				Value []int `validate:"min:10"`
			}{Value: []int{-1, 0, 9, 11}},
			expectedErrCount: 3,
			expectedErr:      ErrBelowMinimum,
		},
		{
			name: "valid 'max' check",
			in: struct {
				Value int `validate:"max:10"`
			}{Value: 2},
		},
		{
			name: "failed 'max' check",
			in: struct {
				Value int `validate:"max:10"`
			}{Value: 31},
			expectedErrCount: 1,
			expectedErr:      ErrAboveMaximum,
		},
		{
			name: "valid 'max' check slice",
			in: struct {
				Value []int `validate:"max:10"`
			}{Value: []int{2, 3, 4}},
		},
		{
			name: "failed 'max' check slice",
			in: struct {
				Value []int `validate:"max:10"`
			}{Value: []int{31, -10}},
			expectedErrCount: 1,
			expectedErr:      ErrAboveMaximum,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, errs)
				return
			}

			var validationErrors ValidationErrors
			require.ErrorAs(t, errs, &validationErrors)
			require.Equal(t, tt.expectedErrCount, len(validationErrors))

			for _, err := range validationErrors {
				require.ErrorIs(t, err.Err, tt.expectedErr)
			}
		})
	}
}
