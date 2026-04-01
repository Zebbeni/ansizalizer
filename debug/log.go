// Package debug provides optional file logging gated behind the -debug flag.
// When disabled (the default), all Log calls are no-ops and no files are created.
package debug

import (
	"fmt"
	"os"
	"time"
)

var (
	logFile *os.File
	enabled bool
)

// Enable turns on debug logging. Call this early in main when the -debug flag is set.
func Enable() {
	enabled = true
}

// Init creates the console.log file if debug mode is enabled.
func Init() {
	if !enabled {
		return
	}
	var err error
	logFile, err = os.Create("console.log")
	if err != nil {
		return
	}
	Log("=== ansizalizer started ===")
}

// Log writes a timestamped message to console.log. No-op if debug is disabled.
func Log(format string, args ...interface{}) {
	if logFile == nil {
		return
	}
	ts := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(logFile, "[%s] %s\n", ts, msg)
	logFile.Sync()
}

// Enabled returns whether debug mode is active.
func Enabled() bool {
	return enabled
}

// Close flushes and closes the log file.
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
