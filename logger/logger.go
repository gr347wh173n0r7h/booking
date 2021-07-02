package logger

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/booking/config"
)

// NewLogger returns a new configured implementation of logrus
func NewLogger(config *config.Config) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(loglevel(config.LogLevel))

	return log
}

func loglevel(l string) logrus.Level {
	level := strings.ToLower(l)
	switch {
	case level == "error":
		return logrus.ErrorLevel
	case level == "info":
		return logrus.InfoLevel
	case level == "debug":
		return logrus.InfoLevel
	case level == "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
