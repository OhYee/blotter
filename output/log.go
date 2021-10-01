package output

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/OhYee/rainbow/color"
	"github.com/OhYee/rainbow/log"
)

var (
	// ErrOutput meesage output
	ErrOutput = log.New().SetColor(color.New().SetFrontRed().SetFontBold()).SetPrefix(withInfo)
	// LogOutput message output
	LogOutput = log.New().SetSuffix(withNewLine)
	// DebugOutput debug message output
	DebugOutput = log.New().SetColor(color.New().SetFrontYellow()).SetSuffix(withNewLine).SetPrefix(withInfo)
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

func withInfo(s string) string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(4, pc)
	thisFile := ""
	callers := make([]string, n)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		if thisFile == "" {
			thisFile = file
		}
		callers[i] = fmt.Sprintf("%s:%d %s", file, line, f.Name())
	}
	return fmt.Sprintf("%s\n\t", strings.Join(callers, "\n"))
}
