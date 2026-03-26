package debug

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func Init() {
	var err error
	logFile, err = os.Create("console.log")
	if err != nil {
		return
	}
	Log("=== ansizalizer started ===")
}

func Log(format string, args ...interface{}) {
	if logFile == nil {
		return
	}
	ts := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(logFile, "[%s] %s\n", ts, msg)
	logFile.Sync()
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
