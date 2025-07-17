package logger

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"mgarnier11.fr/go/libs/utils"
)

const VerboseLevel log.Level = log.DebugLevel - 1
const verboseLevelString = "verbose"

var appLogger *Logger

func InitAppLogger(appName string) *Logger {
	styles := log.DefaultStyles()
	styles.Levels[VerboseLevel] = lipgloss.
		NewStyle().
		SetString(strings.ToUpper(verboseLevelString)).Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("92"))

	log.SetLevel(getLogLevel())
	log.SetStyles(styles)

	appLogger = NewLogger("", "", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

	return appLogger
}

func getLogLevel() log.Level {
	stringLogLevel := utils.GetEnv("LOG_LEVEL", "info")

	level := log.InfoLevel

	if strings.ToLower(stringLogLevel) == verboseLevelString {
		level = VerboseLevel
	} else {
		level, _ = log.ParseLevel(stringLogLevel)
	}

	return level
}

func Debugf(format string, args ...interface{}) {
	appLogger.logf(log.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	appLogger.logf(log.InfoLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	appLogger.logf(log.WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	appLogger.logf(log.ErrorLevel, format, args...)
}

func Verbosef(format string, args ...interface{}) {
	appLogger.logf(VerboseLevel, format, args...)
}
