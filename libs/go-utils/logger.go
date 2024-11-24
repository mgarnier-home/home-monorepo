package goUtils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const VerboseLevel log.Level = log.DebugLevel - 1
const verboseLevelString = "verbose"

func InitLogger() {
	styles := log.DefaultStyles()
	styles.Levels[VerboseLevel] = lipgloss.
		NewStyle().
		SetString(strings.ToUpper(verboseLevelString)).Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("92"))

	stringLogLevel := GetEnv("LOG_LEVEL", "info")

	level := log.InfoLevel

	if strings.ToLower(stringLogLevel) == verboseLevelString {
		level = VerboseLevel
	} else {
		level, _ = log.ParseLevel(stringLogLevel)
	}

	log.SetLevel(level)
	log.SetStyles(styles)
}

func Verbosef(format string, args ...interface{}) {
	log.Logf(VerboseLevel, format, args...)
}
