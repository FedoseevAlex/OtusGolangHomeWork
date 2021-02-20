package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testDataDir = "testdata"
	testEnvDir  = testDataDir + "/env"
)

type testPath struct {
	// Path to create for testing.
	Path string
	// If Path should be directory this field should be true, false if file should be created.
	IsDir bool
	// Data to be written to file.
	// This field will only be used when IsDir is false.
	Data string
}

type readDirTestCase struct {
	// test name
	Name string
	// directory with environment to test if set to empty
	// string then temp dir will be created
	EnvDir string
	// these paths would be created
	// via ioutil.TempFile or ioutil.TempDir under EnvDir
	Create []testPath
	// expected result to compare in test
	Expected Environment
	// expected error
	ExpectedError error
}

var readDirTestCases = []readDirTestCase{
	{
		Name:   "basic functionality",
		EnvDir: testEnvDir,
		Expected: Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		},
	},
	{
		Name: "check skip subdirectories",
		Create: []testPath{
			{Path: "dir_to_be_skipped", IsDir: true},
			{Path: "VARIABLE", Data: "value"},
		},
		Expected: Environment{
			"VARIABLE": {Value: "value", NeedRemove: false},
		},
	},
	{
		Name: "check error for variables with = in name",
		Create: []testPath{
			{Path: "VARIABLE", Data: "value"},
			{Path: "VARIABLE=", Data: "This won't be in env"},
		},
		ExpectedError: ErrWrongFileName,
	},
	{
		Name: "check error for variables with whitespace in name",
		Create: []testPath{
			{Path: "VARIABLE", Data: "value"},
			{Path: "VARIABLE WITH SPACES", Data: "This value will be ignored too"},
		},
		ExpectedError: ErrWrongFileName,
	},
}

func TestReadDir(t *testing.T) {
	for _, testCase := range readDirTestCases {
		data := testCase
		t.Run(data.Name, func(t *testing.T) {
			setUpReadDirCase(t, &data)

			env, err := ReadDir(data.EnvDir)
			require.ErrorIs(t, err, data.ExpectedError)
			require.Equal(t, env, data.Expected)
		})
	}
}

func setUpReadDirCase(t *testing.T, testCase *readDirTestCase) {
	if _, err := os.Stat(testCase.EnvDir); os.IsNotExist(err) {
		testCase.EnvDir = t.TempDir()
	}

	for _, path := range testCase.Create {
		if path.IsDir {
			err := os.MkdirAll(filepath.Join(testCase.EnvDir, path.Path), 0o664)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			file, err := os.Create(filepath.Join(testCase.EnvDir, path.Path))
			if err != nil {
				t.Fatal(err)
			}

			_, err = file.WriteString(path.Data)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

type extractValueTestCase struct {
	name     string
	contents []byte
	expected string
}

var extractValueTestCases = []extractValueTestCase{
	{
		name:     "basic case",
		contents: []byte("example value"),
		expected: "example value",
	},
	{
		name:     "case with zero byte",
		contents: []byte("line\x00feed"),
		expected: "line\nfeed",
	},
	{
		name:     "ensure trim trailing tabs",
		contents: []byte("\tsome\ttabs\t\t\t\t\t\t\t\t\t\t\t"),
		expected: "\tsome\ttabs",
	},
	{
		name:     "ensure trim trailing spaces",
		contents: []byte(" some spaces     "),
		expected: " some spaces",
	},
	{
		name:     "ensure trim trailing tabs and spaces",
		contents: []byte("tab space mix\t \t \t"),
		expected: "tab space mix",
	},
}

func TestExtractValue(t *testing.T) {
	for _, testCase := range extractValueTestCases {
		data := testCase
		t.Run(data.name, func(t *testing.T) {
			res := extractValue(data.contents)
			require.Equal(t, res, data.expected)
		})
	}
}
