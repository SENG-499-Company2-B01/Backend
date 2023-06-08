package logger_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"fmt"

	"github.com/SENG-499-Company2-B01/Backend/logger"
)

func TestLoggerInfo(t *testing.T) {
	// Redirect log output to a buffer for testing
	var buf bytes.Buffer

	// Initialize the logger with the buffer as the output
	logger.InitLogger(&buf, &buf, &buf)

	// Test Info
	expectedInfo := "\033[1;34m[INFO]\033[0m TestInfoMessage"
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
	logger.InitLogger(&buf, &buf, &buf)

	// Test Warning
	expectedWarning := "\033[1;33m[WARNING]\033[0m logger_test.go:43 TestWarningMessage"
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
	logger.InitLogger(&buf, &buf, &buf)

	// Turn off the exit on error to test successfully
	logger.SetExitOnError(false)

	// Test Error
	expectedError := "\033[1;31m[ERROR]\033[0m logger_test.go:66 TestErrorMessage"
	logger.Error(fmt.Errorf("TestErrorMessage"))

	// Turn back on the exit on error 
	logger.SetExitOnError(true)

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

	// Get the last line of the log output
	lastLine := lines[len(lines)-1]

	// Split the last line by space characters
	// The third element should be the log message without the date and time
	parts := strings.SplitN(lastLine, " ", 3)

	// Get the log message without the date and time
	lastLineWithoutTimestamp := parts[2]

	return lastLineWithoutTimestamp
}

func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Clean up
	log.SetOutput(os.Stderr)

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}