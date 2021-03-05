package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type runCmdTestCase struct {
	Name         string
	Command      []string
	ExpectedCode int
}

var runCmdTestCases = []runCmdTestCase{
	{
		Name:         "check 0 exit code",
		Command:      []string{"/bin/bash", "-c", "true"},
		ExpectedCode: 0,
	},
	{
		Name:         "check error exit code",
		Command:      []string{"/bin/bash", "-c", "false"},
		ExpectedCode: 1,
	},
	{
		Name:         "check for envdir exit code",
		Command:      []string{"command_not_exist"},
		ExpectedCode: 111,
	},
}

func TestRunCmd(t *testing.T) {
	testEnvironment := Environment{}

	for _, testCase := range runCmdTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			returnCode := RunCmd(testCase.Command, testEnvironment)
			require.Equal(t, testCase.ExpectedCode, returnCode)
		})
	}
}

type prepareEnvTestCase struct {
	Name     string
	Env      Environment
	Expected []string
	// Here will be names of variables that must be
	// abscent after prepareEnv.
	NotExpected []string
}

var prepareEnvTestCases = []prepareEnvTestCase{
	{
		Name: "basic functionality",
		Env: Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		},
		Expected: []string{
			"BAR=bar",
			"EMPTY=",
			"FOO=   foo\nwith new line",
			"HELLO=\"hello\"",
		},
		NotExpected: []string{
			"UNSET",
		},
	},
}

func TestPrepareEnv(t *testing.T) {
	for _, testCase := range prepareEnvTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := prepareEnv(testCase.Env)
			require.NoError(t, err)
			require.Subset(t, os.Environ(), testCase.Expected)

			for _, variable := range testCase.NotExpected {
				_, varIsSet := os.LookupEnv(variable)
				require.False(t, varIsSet)
			}
		})
	}
}
