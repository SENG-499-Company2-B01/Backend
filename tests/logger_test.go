package tests

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SENG-499-Company2-B01/Backend/logger"
)

func TestLoggerInfo(t *testing.T) {
	// Redirect log output to a buffer for testing
	var buf bytes.Buffer

	// Initialize the logger with the buffer as the output
	logger.InitLogger(&buf, &buf, &buf, nil)

	// Test Info
	expectedInfo := fmt.Sprintf("%s \033[1;34m[INFO]\033[0m      TestInfoMessage", getTimestamp())
	logger.Info("TestInfoMessage")

	// Verify Info log message
	if got := extractLogMessage(buf.String()); got != expectedInfo {
		t.Errorf("Info log message:\nExpected: %q\nGot: %q", expectedInfo, got)
	}

	// Clean up
	log.SetOutput(os.Stderr)
}

func TestLoggerWarning(t *testing.T) {
	// Redirect log output to a buffer for testing
	var buf bytes.Buffer

	// Initialize the logger with the buffer as the output
	logger.InitLogger(&buf, &buf, &buf, nil)

	// Test Warning
	expectedWarning := fmt.Sprintf("%s \033[1;33m[WARNING]\033[0m   logger_test.go:44 | TestWarningMessage", getTimestamp())
	logger.Warning("TestWarningMessage")

	// Verify Warning log message
	if got := extractLogMessage(buf.String()); got != expectedWarning {
		t.Errorf("Warning log message:\nExpected: %q\nGot: %q", expectedWarning, got)
	}

	// Clean up
	log.SetOutput(os.Stderr)
}

func TestLoggerError(t *testing.T) {
	// Redirect log output to a buffer for testing
	var buf bytes.Buffer

	// Initialize the logger with the buffer as the output
	logger.InitLogger(&buf, &buf, &buf, nil)

	// Test Error
	expectedError := fmt.Sprintf("%s \033[1;31m[ERROR]\033[0m     logger_test.go:64 | TestErrorMessage", getTimestamp())
	logger.Error(fmt.Errorf("TestErrorMessage"))

	// Verify Error log message
	if got := extractLogMessage(buf.String()); got != expectedError {
		t.Errorf("Error log message:\nExpected: %q\nGot: %q", expectedError, got)
	}

	// Clean up
	log.SetOutput(os.Stderr)
}

func extractLogMessage(logOutput string) string {
	// Split log output by newline characters
	lines := strings.Split(strings.TrimSpace(logOutput), "\n")

	// The log message doesn't contain file info, return the entire last line
	return lines[len(lines)-1]
}

func getTimestamp() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Clean up
	log.SetOutput(os.Stderr)

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
