package hw09_struct_validator //nolint:golint,stylecheck,revive

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIntInValidator(t *testing.T) {
	tests := []struct {
		name        string
		inputTag    string
		result      []int
		expectedErr error
	}{
		{
			name:     "valid case",
			inputTag: "in:100,200,300",
			result:   []int{100, 200, 300},
		},
		{
			name:     "valid case: negative numbers",
			inputTag: "in:100,-200,-300",
			result:   []int{100, -200, -300},
		},
		{
			name:        "not a number choice",
			inputTag:    "in:100,wow,hey,stop",
			expectedErr: ErrTagParse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseIntInValidatorTag(tt.inputTag)

			require.Equal(t, tt.result, result)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestParseIntMaxValidator(t *testing.T) {
	tests := []struct {
		name        string
		inputTag    string
		result      int
		expectedErr error
	}{
		{
			name:     "valid case",
			inputTag: "max:100",
			result:   100,
		},
		{
			name:     "valid case: negative number",
			inputTag: "max:-100",
			result:   -100,
		},
		{
			name:        "not a number",
			inputTag:    "max:wow",
			expectedErr: ErrTagParse,
		},
		{
			name:        "corrupted tag",
			inputTag:    "max:",
			expectedErr: ErrCorruptedTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseIntMaxValidatorTag(tt.inputTag)

			require.Equal(t, tt.result, result)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestParseIntMinValidator(t *testing.T) {
	tests := []struct {
		name        string
		inputTag    string
		result      int
		expectedErr error
	}{
		{
			name:     "valid case",
			inputTag: "min:100",
			result:   100,
		},
		{
			name:     "valid case: negative number",
			inputTag: "min:-100",
			result:   -100,
		},
		{
			name:        "not a number",
			inputTag:    "min:wow",
			expectedErr: ErrTagParse,
		},
		{
			name:        "corrupted tag",
			inputTag:    "min:",
			expectedErr: ErrCorruptedTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseIntMaxValidatorTag(tt.inputTag)

			require.Equal(t, tt.result, result)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestParseStrLenValidator(t *testing.T) {
	tests := []struct {
		name        string
		inputTag    string
		result      int
		expectedErr error
	}{
		{
			name:     "valid case",
			inputTag: "len:100",
			result:   100,
		},
		{
			name:        "valid case: negative number",
			inputTag:    "len:-100",
			expectedErr: ErrTagParse,
		},
		{
			name:        "not a number",
			inputTag:    "len:wow",
			expectedErr: ErrTagParse,
		},
		{
			name:        "corrupted tag",
			inputTag:    "len:",
			expectedErr: ErrCorruptedTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseStrLenValidatorTag(tt.inputTag)

			require.Equal(t, tt.result, result)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestParseRegexpValidator(t *testing.T) {
	tests := []struct {
		name        string
		inputTag    string
		result      *regexp.Regexp
		expectedErr error
	}{
		{
			name:     "valid case",
			inputTag: "regexp:\\d",
			result:   regexp.MustCompile(`\d`),
		},
		{
			name:        "wrong regexp",
			inputTag:    "regexp:\\d(",
			expectedErr: ErrTagParse,
		},
		{
			name:        "corrupted tag",
			inputTag:    "regexp:",
			expectedErr: ErrCorruptedTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRegexpValidatorTag(tt.inputTag)

			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.result, result)
		})
	}
}

func TestParseStrInValidator(t *testing.T) {
	tests := []struct {
		name        string
		inputTag    string
		result      []string
		expectedErr error
	}{
		{
			name:     "valid case",
			inputTag: "in:foo,bar,a,b",
			result:   []string{"foo", "bar", "a", "b"},
		},
		{
			name:        "corrupted tag",
			inputTag:    "in:",
			expectedErr: ErrCorruptedTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseStrInValidatorTag(tt.inputTag)

			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.result, result)
		})
	}
}
