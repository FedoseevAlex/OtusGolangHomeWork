package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testDir   = "testdata/"
	inputPath = testDir + "input.txt"
)

type copyTestCase struct {
	// Path to file with input data
	fromPath string
	// Path with correct copy results
	referencePath string

	offset int64
	limit  int64
}

var copyTestCases = []copyTestCase{
	{
		fromPath:      inputPath,
		referencePath: testDir + "out_offset0_limit0.txt",
		offset:        0,
		limit:         0,
	},
	{
		fromPath:      inputPath,
		referencePath: testDir + "out_offset0_limit10.txt",
		offset:        0,
		limit:         10,
	},
	{
		fromPath:      inputPath,
		referencePath: testDir + "out_offset0_limit1000.txt",
		offset:        0,
		limit:         1000,
	},
	{
		fromPath:      inputPath,
		referencePath: testDir + "out_offset0_limit10000.txt",
		offset:        0,
		limit:         10000,
	},
	{
		fromPath:      inputPath,
		referencePath: testDir + "out_offset100_limit1000.txt",
		offset:        100,
		limit:         1000,
	},
	{
		fromPath:      inputPath,
		referencePath: testDir + "out_offset6000_limit1000.txt",
		offset:        6000,
		limit:         1000,
	},
}

func TestCopyBasic(t *testing.T) {

	for _, testData := range copyTestCases {
		data := testData
		name := fmt.Sprintf("test offset:%d limit:%d", data.offset, data.limit)

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tempFileNamePattern := fmt.Sprintf(
				"offset%d_limit%d.result.*.txt",
				data.offset,
				data.limit,
			)
			tempFile, err := ioutil.TempFile(testDir, tempFileNamePattern)
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			err = Copy(data.fromPath, tempFile.Name(), data.offset, data.limit)
			require.NoError(t, err)

			// compare out file and referencePath contents
			expected, err := ioutil.ReadFile(data.referencePath)
			require.NoError(t, err)

			result, err := ioutil.ReadFile(tempFile.Name())
			require.NoError(t, err)

			require.True(t, bytes.Equal(expected, result))
		})
	}
}

type validateTestCase struct {
	// test case name
	name string

	// input parameters
	from   string
	to     string
	offset int64
	limit  int64

	// expected output
	expected error
}

var validateTestCases = []validateTestCase{
	{
		name:     "offset is larger than file size",
		from:     inputPath,
		to:       "output/file/path",
		offset:   10000,
		limit:    0,
		expected: ErrOffsetExceedsFileSize,
	},
	{
		name:     "valid case",
		from:     inputPath,
		to:       "output/file/path",
		offset:   0,
		limit:    0,
		expected: nil,
	},
	{
		name:     "negative offset",
		from:     inputPath,
		to:       "output/file/path",
		offset:   -10,
		limit:    0,
		expected: ErrOffsetIsNegative,
	},
	{
		name:     "negative limit",
		from:     inputPath,
		to:       "output/file/path",
		offset:   0,
		limit:    -239,
		expected: ErrLimitIsNegative,
	},
	{
		name:     "attempt to copy endless file",
		from:     "/dev/urandom",
		to:       "output/file/path",
		offset:   0,
		limit:    0,
		expected: ErrReadWithoutLimit,
	},
	{
		name:     "no src file specified",
		to:       "output/file/path",
		offset:   0,
		limit:    0,
		expected: ErrNoSrcFileSpecified,
	},
	{
		name:     "no dst file specified",
		from:     inputPath,
		offset:   0,
		limit:    0,
		expected: ErrNoDstFileSpecified,
	},
}

func TestValidate(t *testing.T) {
	for _, testData := range validateTestCases {
		data := testData
		t.Run(data.name, func(t *testing.T) {
			err := Validate(data.from, data.to, data.offset, data.limit)
			require.ErrorIs(t, err, data.expected)
		})
	}
}
