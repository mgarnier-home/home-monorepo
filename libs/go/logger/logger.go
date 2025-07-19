package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	name       string
	nameFormat string
	style      lipgloss.Style

	parent *Logger
}

func (logger *Logger) Debug(msg string) {
	logger.logf(log.DebugLevel, msg, nil)
}

func (logger *Logger) Info(msg string) {
	logger.logf(log.InfoLevel, msg, nil)
}

func (logger *Logger) Warn(msg string) {
	logger.logf(log.WarnLevel, msg, nil)
}

func (logger *Logger) Error(msg string) {
	logger.logf(log.ErrorLevel, msg, nil)
}

func (logger *Logger) Verbose(msg string) {
	logger.logf(VerboseLevel, msg, nil)
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
	var nameString string

	if len(logger.name) > 0 {
		nameString = logger.style.Render(fmt.Sprintf(logger.nameFormat, logger.name))
	}

	coloredString := logger.style.Render(fmt.Sprintf(format, args...))

	if logger.parent != nil {
		logger.parent.logf(level, nameString+coloredString)
	} else {
		log.Logf(level, nameString+coloredString)
	}

}

func NewLogger(name string, nameFormat string, style lipgloss.Style, parent *Logger) *Logger {
	if parent == nil {
		parent = appLogger
	}

	return &Logger{
		name:       name,
		nameFormat: nameFormat,
		style:      style,
		parent:     parent,
	}
}
