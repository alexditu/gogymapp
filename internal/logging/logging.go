package logging

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"fmt"
)

const (
	defaultTimestampFormat = "2023-12-31 16:00:00.999"
)

type Settings struct {
	LogToFile       bool
	FileName        string
	TimestampFormat string
	Level           log.Level
}

func InitStandardLogger(setts *Settings) error {
	if setts == nil {
		return fmt.Errorf("nil settings")
	}

	return initLogger(log.StandardLogger(), setts)
}

func NewLogger(setts *Settings) (*log.Logger, error) {
	l := log.New()
	if l == nil {
		return nil, fmt.Errorf("failed to create new logger object")
	}

	if setts == nil {
		return nil, fmt.Errorf("nil settings")
	}

	err := initLogger(l, setts)

	return l, err
}

func initLogger(logger *log.Logger, setts *Settings) error {
	if setts.TimestampFormat == "" {
		setts.TimestampFormat = defaultTimestampFormat
	}

	logger.SetLevel(setts.Level)

	if setts.LogToFile {
		if setts.FileName == "" {
			return fmt.Errorf("invalid settings: logToFile is true but fileName is empty")
		}

		initFileLogger(logger, setts.FileName)
	}

	logger.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		TimestampFormat: setts.TimestampFormat,
		FullTimestamp:   true,
	})

	return nil
}

func initFileLogger(logger *log.Logger, fileName string) {
	logger.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    50, // MiB
		MaxBackups: 3,
		MaxAge:     60,   //days
		Compress:   true, // disabled by default
	})
}
