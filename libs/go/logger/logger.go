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
