package main

import (
	"log"
	"os"
)

func main() {
	envDir := os.Args[1]
	cmdAndArgs := os.Args[2:]

	environment, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(cmdAndArgs, environment)
	os.Exit(returnCode)
}
