package output

import (
	"github.com/OhYee/rainbow/color"
	"github.com/OhYee/rainbow/log"
)

var (
	// ErrOutput meesage output
	ErrOutput = log.New().SetColor(color.New().SetFrontRed().SetFontBold())
	// LogOutput message output
	LogOutput = log.New().SetSuffix(withNewLine)
	// DebugOutput debug message output
	DebugOutput = log.New().SetColor(color.New().SetFrontYellow()).SetSuffix(withNewLine)
)

// Err meage output
func Err(err error) {
	ErrOutput.Println(err.Error())
}

// Log message output
func Log(format string, args ...interface{}) {
	LogOutput.Printf(format, args...)
}

// Debug message output
func Debug(format string, args ...interface{}) {
	DebugOutput.Printf(format, args...)
}

func withNewLine(s string) string {
	return "\n"
}
