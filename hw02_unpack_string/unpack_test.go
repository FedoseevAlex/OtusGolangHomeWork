package hw02unpackstring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for i, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "те5ст",
			expected: "тееееест",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa13b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
	} {
		result, err := Unpack(tst.input)
		fmt.Printf("Test iteration: %d | res: %s expect: %s\n", i, result, tst.expected)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for i, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `qwe\`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `qw\xe`,
			expected: "",
			err:      ErrInvalidString,
		},
	} {
		result, err := Unpack(tst.input)
		fmt.Printf("Test iteration: %d | res: %s expect: %s\n", i, result, tst.expected)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
