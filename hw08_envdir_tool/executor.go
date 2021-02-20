package main

import (
	"log"
	"os"
	"os/exec"
)

// Just like the original envdir.
const envDirErrorReturnCode = 111

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Cmd{}
	command.Path = cmd[0]
	command.Args = append(command.Args, cmd...)

	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	err := prepareEnv(env)
	if err != nil {
		log.Println(err)
		returnCode = envDirErrorReturnCode
	}

	// Check if err IS NOT nil and IS NOT of type *exec.ExitError. Set returnCode to 111 in that case.
	// Otherwise set returnCode to command run exit code.
	err = command.Run()
	// Here we skip linting because it is necessary to check for error type.
	if _, ok := err.(*exec.ExitError); !ok && err != nil { //nolint:errorlint
		log.Println(err)
		returnCode = envDirErrorReturnCode
	} else {
		returnCode = command.ProcessState.ExitCode()
	}

	return
}

func prepareEnv(env Environment) error {
	for varName, varValue := range env {
		if varValue.NeedRemove {
			err := os.Unsetenv(varName)
			if err != nil {
				return err
			}
			continue
		}
		err := os.Setenv(varName, varValue.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
