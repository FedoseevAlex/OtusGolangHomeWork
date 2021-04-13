package logger

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	testDir := t.TempDir()
	logfile, err := ioutil.TempFile(testDir, "logger_test_*")
	require.NoError(t, err)
	log := New("info", logfile.Name())

	t.Run("simple log entry", func(t *testing.T) {
		args := make(LogArgs)
		args["error1"] = errors.New("some error")
		args["host"] = "localhost"
		log.Info("Test message", args)

		logs, err := ioutil.ReadAll(logfile)
		stringLogs := string(logs)

		require.NoError(t, err)
		require.Contains(t, stringLogs, "Test message")
		require.Contains(t, stringLogs, "\"error1\":\"some error\"")
	})

	t.Run("skip entry by level", func(t *testing.T) {
		log.Debug("THIS MUST NOT BE IN LOG")

		logs, err := ioutil.ReadAll(logfile)
		stringLogs := string(logs)

		require.NoError(t, err)
		require.NotContains(t, stringLogs, "THIS MUST NOT BE IN LOG")
	})
}
