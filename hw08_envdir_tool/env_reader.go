package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrWrongFileName = errors.New("wrong file name")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		return walkFunc(env, dir, path, info, err)
	})
	if err != nil {
		return nil, err
	}
	return env, nil
}

func walkFunc(env Environment, dir, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	ok := validateFileName(path)
	if !ok {
		return ErrWrongFileName
	}

	// Skip all subdirectories of dir
	if info.IsDir() {
		if path == dir {
			return nil
		}
		return filepath.SkipDir
	}

	varName, varValue, err := processFile(path)
	if err != nil {
		return err
	}

	env[varName] = varValue
	return nil
}

func validateFileName(path string) bool {
	return !strings.ContainsAny(path, " =")
}

func processFile(path string) (varName string, varValue EnvValue, err error) {
	varName = filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return
	}

	if info.Size() == 0 {
		varValue.NeedRemove = true
		return
	}

	// scan only one line
	scan := bufio.NewScanner(file)

	if ok := scan.Scan(); ok {
		varValue.Value = extractValue(scan.Bytes())
	}
	return
}

func extractValue(data []byte) string {
	replaced := bytes.ReplaceAll(data, []byte{0}, []byte{'\n'})
	replaced = bytes.TrimRight(replaced, " \t")
	return string(replaced)
}
