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
		StrRegexpSlice []string `validate:"regexp:^a+$"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "valid length case",
			in:          StrLenCheck{StrLen: "five!"},
			expectedErr: nil,
		}, {
			name:        "invalid length case",
			in:          StrLenCheck{StrLen: "length is not five"},
			expectedErr: ErrStrLengthInvalid,
		},
		{
			name:        "valid length case in slice",
			in:          StrLenSliceCheck{StrLenSlice: []string{"aa", "bb"}},
			expectedErr: nil,
		},
		{
			name:        "invalid length case in slice",
			in:          StrLenSliceCheck{StrLenSlice: []string{"a", "bb"}},
			expectedErr: ErrStrLengthInvalid,
		},
		{
			name:        "invalid length case in slice",
			in:          StrLenSliceCheck{StrLenSlice: []string{"aaa", "bb"}},
			expectedErr: ErrStrLengthInvalid,
		},
		{
			name:        "valid regexp check",
			in:          StrRegexpCheck{StrRegexp: "aa"},
			expectedErr: nil,
		},
		{
			name:        "invalid regexp check",
			in:          StrRegexpCheck{StrRegexp: "adfadfa"},
			expectedErr: ErrRegexpMismatch,
		},
		{
			name:        "valid regexp check in slice",
			in:          StrRegexpSliceCheck{StrRegexpSlice: []string{"aa", "aaaaaa"}},
			expectedErr: nil,
		},
		{
			name:        "invalid regexp check in slice",
			in:          StrRegexpSliceCheck{StrRegexpSlice: []string{"aa", "aaaaba", ""}},
			expectedErr: ErrRegexpMismatch,
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

			for _, err := range validationErrors {
				require.ErrorIs(t, err.Err, tt.expectedErr)
			}
		})
	}
}
