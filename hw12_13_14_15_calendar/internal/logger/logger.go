package logger

import (
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type (
	LogArgs map[string]interface{}
	Logger  struct {
		logger zerolog.Logger
	}
)

func New(level, file string) *Logger {
	logfile, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o640)
	if err != nil {
		log.Println("can't open file ", file)
		return nil
	}

	logger := zerolog.New(logfile).With().Timestamp().Logger()
	numLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		logger.Warn().
			Str("level", level).
			Msg("Unknown log level. Using DEBUG.")
	} else {
		logger = logger.Level(numLevel)
	}

	return &Logger{logger: logger}
}

func (l Logger) Trace(msg string, args ...map[string]interface{}) {
	l.logWithLevel(zerolog.TraceLevel, msg, args...)
}

func (l Logger) Debug(msg string, args ...map[string]interface{}) {
	l.logWithLevel(zerolog.DebugLevel, msg, args...)
}

func (l Logger) Info(msg string, args ...map[string]interface{}) {
	l.logWithLevel(zerolog.InfoLevel, msg, args...)
}

func (l Logger) Warn(msg string, args ...map[string]interface{}) {
	l.logWithLevel(zerolog.WarnLevel, msg, args...)
}

func (l Logger) Error(msg string, args ...map[string]interface{}) {
	l.logWithLevel(zerolog.ErrorLevel, msg, args...)
}

func (l Logger) logWithLevel(level zerolog.Level, msg string, args ...map[string]interface{}) {
	record := l.logger.WithLevel(level)

	for _, argSet := range args {
		record = record.Fields(argSet)
	}
	record.Msg(msg)
}
