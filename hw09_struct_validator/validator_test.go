package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"encoding/json"
	"fmt"
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
		Version  string   `validate:"len:5"`
		Features []string `validate:"len:2"`
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{Version: "five!", Features: []string{"aa", "bb", "cc"}},
			expectedErr: nil,
		},
		{
			in:          App{Version: "length is not five", Features: []string{"aa", "bb", "cc", "dd"}},
			expectedErr: ErrStrLengthInvalid,
		},
		{
			in:          App{Version: "five!", Features: []string{"aaa", "bb", "cc", "ddd"}},
			expectedErr: ErrStrLengthInvalid,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
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
