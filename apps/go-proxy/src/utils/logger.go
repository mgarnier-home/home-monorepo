package utils

import (
	"goUtils"

	"github.com/charmbracelet/log"
)

type Logger struct {
	name       string
	nameFormat string

	parent *Logger
}

func NewLogger(name string, nameFormat string, parent *Logger) *Logger {
	return &Logger{
		name:       name,
		nameFormat: nameFormat,
		parent:     parent,
	}
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logf(log.DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logf(log.InfoLevel, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logf(log.WarnLevel, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logf(log.ErrorLevel, format, args...)
}

func (logger *Logger) Verbosef(format string, args ...interface{}) {
	logger.logf(goUtils.VerboseLevel, format, args...)
}

func (logger *Logger) logf(level log.Level, format string, args ...interface{}) {
	if len(args) > 0 {
		args = append([]interface{}{logger.name}, args...)
	} else {
		args = []interface{}{logger.name}
	}

	if logger.parent != nil {
		logger.parent.logf(level, logger.nameFormat+format, args...)
	} else {
		log.Logf(level, logger.nameFormat+format, args...)

	}

}
