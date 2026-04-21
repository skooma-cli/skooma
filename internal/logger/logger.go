// Internal package for logging
package logger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

var l *log.Logger

// Init initializes the logger.
func Init() error {
	logPath, err := GetLogFilePath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l = log.NewWithOptions(f, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		ReportCaller:    true,
		CallerOffset:    1,
		Level:           log.DebugLevel,
	})

	return nil
}

// Info logs an info message.
func Info(msg string, keyvals ...any) { l.Info(msg, keyvals...) }

// Debug logs a debug message.
func Debug(msg string, keyvals ...any) { l.Debug(msg, keyvals...) }

// Warn logs a warning message.
func Warn(msg string, keyvals ...any) { l.Warn(msg, keyvals...) }

// Error logs an error message.
func Error(msg string, keyvals ...any) { l.Error(msg, keyvals...) }

// Fatal logs a fatal message and exits.
func Fatal(msg string, keyvals ...any) {
	fmt.Printf("❌ %s\n\n%s\n\nRun `skooma log` for more details\n", msg, strings.Repeat("-", 80))
	l.Fatal(msg, keyvals...)
}

// GetLogFilePath returns the path to the log file.
func GetLogFilePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, "skooma", "skooma.log"), nil
}

// ViewLog opens the log file in the user's default pager.
func ViewLog() error {
	logPath, err := GetLogFilePath()
	if err != nil {
		return err
	}

	pager := os.Getenv("PAGER")
	if pager == "" {
		switch runtime.GOOS {
		case "windows":
			pager = "more"
		default:
			pager = "less"
		}
	}

	file, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command(pager)
	cmd.Stdin = file
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
