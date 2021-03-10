package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}
	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	// Structs to test every validator separately
	StrLenCheck struct {
		StrLen string `validate:"len:5"`
	}
	StrLenSliceCheck struct {
		StrLenSlice []string `validate:"len:2"`
	}

	StrRegexpCheck struct {
		StrRegexp string `validate:"regexp:^a+$"`
	}
	StrRegexpSliceCheck struct {
		StrRegexpSlice []string `validate:"regexp:^\\d+$"`
	}

	StrInCheck struct {
		StrIn string `validate:"in:a,b,c,d"`
	}
	StrInSliceCheck struct {
		StrInSlice []string `validate:"in:a,b,c,d"`
	}
)

//func TestValidate(t *testing.T) {
//	tests := []struct {
//		name        string
//		in          interface{}
//		expectedErr error
//	}{}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			errs := Validate(tt.in)
//
//			if tt.expectedErr == nil {
//				require.NoError(t, errs)
//				return
//			}
//
//			var validationErrors ValidationErrors
//			require.ErrorAs(t, errs, &validationErrors)
//
//			for _, err := range validationErrors {
//				require.ErrorIs(t, err.Err, tt.expectedErr)
//			}
//		})
//	}
//}

// Here we check basic functionality of every validator.
// So we are interested only in validation errors leaving aside
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
			in:   StrLenCheck{StrLen: "five!"},
		}, {
			name:             "invalid length case",
			in:               StrLenCheck{StrLen: "length is not five"},
			expectedErrCount: 1,
			expectedErr:      ErrStrLengthInvalid,
		},
		{
			name: "valid length case in slice",
			in:   StrLenSliceCheck{StrLenSlice: []string{"aa", "bb"}},
		},
		{
			name:             "invalid length case in slice",
			in:               StrLenSliceCheck{StrLenSlice: []string{"a", "bb", "ccc"}},
			expectedErrCount: 2,
			expectedErr:      ErrStrLengthInvalid,
		},
		{
			name: "valid regexp check",
			in:   StrRegexpCheck{StrRegexp: "aa"},
		},
		{
			name:             "invalid regexp check",
			in:               StrRegexpCheck{StrRegexp: "adfadfa"},
			expectedErrCount: 1,
			expectedErr:      ErrRegexpMismatch,
		},
		{
			name: "valid regexp check in slice",
			in:   StrRegexpSliceCheck{StrRegexpSlice: []string{"12", "234523452345"}},
		},
		{
			name:             "invalid regexp check in slice",
			in:               StrRegexpSliceCheck{StrRegexpSlice: []string{"123", "0xDEADFACECAFE", ""}},
			expectedErrCount: 2,
			expectedErr:      ErrRegexpMismatch,
		},

		{
			name: "valid in check",
			in:   StrInCheck{StrIn: "a"},
		},
		{
			name:             "failed in check",
			in:               StrInCheck{StrIn: "e"},
			expectedErrCount: 1,
			expectedErr:      ErrNotInChoices,
		},
		{
			name: "valid in check in slice",
			in:   StrInSliceCheck{StrInSlice: []string{"a", "b", "c"}},
		},
		{
			name:             "failed in check in slice",
			in:               StrInSliceCheck{StrInSlice: []string{"e", "Not in available choices", "a"}},
			expectedErrCount: 2,
			expectedErr:      ErrNotInChoices,
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
