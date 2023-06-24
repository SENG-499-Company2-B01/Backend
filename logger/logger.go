package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	// Flag to control program termination on error log
	exitOnError = false // Default false

	fileWriter io.Writer
)

type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Initializes the loggers with their respective outputs and formats.
func InitLogger(infoOutput, warningOutput, errorOutput, fileOutput io.Writer) error {
	infoLogger = log.New(infoOutput, "", log.Ldate|log.Ltime)
	warningLogger = log.New(warningOutput, "", log.Ldate|log.Ltime)
	errorLogger = log.New(errorOutput, "", log.Ldate|log.Ltime)

	fileWriter = fileOutput

	// Check if any of the loggers failed to initialize
	if infoLogger == nil || warningLogger == nil || errorLogger == nil {
		return errors.New("Error initializing logger")
	}

	return nil
}

// Returns the file and line number of the calling function.
func getFileLine() string {
	_, file, line, _ := runtime.Caller(2) // Skip two frames: getFileLine() and the calling function
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// SetExitOnError sets the behavior of program termination on error logs.
// If exit is true, the program will exit on error logs.
// If exit is false, the program will not exit on error logs.
func SetExitOnError(exit bool) {
	exitOnError = exit
}

// Info logs an information message to the console with the file and line number.
func Info(message string) {
	// Print the information message to the console and colour the tag blue
	infoLogger.Println("\033[1;34m[INFO]\033[0m     ", message)
}

// Warning logs a warning message to the console with the file and line number.
func Warning(message string) {
	fileLine := getFileLine()

	// Print the warning message to the console and colour the tag yellow
	warningLogger.Println("\033[1;33m[WARNING]\033[0m  ", fileLine, "|", message)

	// Write the warning message to the file
	if fileWriter != nil {
		fileWriter.Write([]byte(fmt.Sprintf("%s [WARNING]   %s | %s\n", getTimestamp(), fileLine, message)))
	}
}

func getErrorMessage(err error, code int) string {
	var errorMessage errorMessage

	errorMessage.Code = code
	errorMessage.Message = err.Error()

	result, _ := json.Marshal(errorMessage)
	return string(result)
}

// Error logs an error message to the console with the file and line number.
// If exitOnError flag is set, the program will exit after logging the error.
func Error(err error, code int) {
	fileLine := getFileLine()

	// Print the error to the console with the file and line number, and colour the tag red
	errorLogger.Printf("\033[1;31m[ERROR]\033[0m     %s | %s\n", fileLine, getErrorMessage(err, code))

	// Write the error to the file
	if fileWriter != nil {
		fileWriter.Write([]byte(fmt.Sprintf("%s [ERROR]     %s | %s\n", getTimestamp(), fileLine, getErrorMessage(err, code))))
	}

	if exitOnError {
		os.Exit(1)
	}
}

// Returns the current timestamp in a specific format.
func getTimestamp() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
