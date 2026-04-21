// Internal package for logging
package logger

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/skooma-cli/skooma/internal/config"
)

var l *log.Logger

func Init() error {
	logPath, err := config.GetLogFilePath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l = log.NewWithOptions(f, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
		Level:           log.DebugLevel,
	})

	return nil
}

func Info(msg string, keyvals ...any)  { l.Info(msg, keyvals...) }
func Debug(msg string, keyvals ...any) { l.Debug(msg, keyvals...) }
func Warn(msg string, keyvals ...any)  { l.Warn(msg, keyvals...) }
func Error(msg string, keyvals ...any) { l.Error(msg, keyvals...) }
func Fatal(msg string, keyvals ...any) { l.Fatal(msg, keyvals...) }
