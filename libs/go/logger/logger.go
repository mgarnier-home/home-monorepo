package logger

import "github.com/charmbracelet/log"

type Logger struct {
	name       string
	nameFormat string

	parent *Logger
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
	logger.logf(VerboseLevel, format, args...)
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

func NewLogger(name string, nameFormat string, parent *Logger) *Logger {
	if parent == nil {
		parent = appLogger
	}

	return &Logger{
		name:       name,
		nameFormat: nameFormat,
		parent:     parent,
	}
}
