package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

var levelColors = map[LogLevel]string{
	DEBUG: "\033[36m", // Cyan
	INFO:  "\033[32m", // Green
	WARN:  "\033[33m", // Yellow
	ERROR: "\033[31m", // Red
	FATAL: "\033[35m", // Magenta
}

const (
	colorReset = "\033[0m"
)

type Logger struct {
	minLevel  LogLevel
	isColored bool
	logger    *log.Logger
}

var defaultLogger *Logger

// ParseLevel converts a string level to LogLevel
func ParseLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return INFO, fmt.Errorf("invalid log level: %s", level)
	}
}

// InitLogger initializes the default logger
func InitLogger(level LogLevel, useColors bool) {
	defaultLogger = NewLogger(level, useColors)
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, useColors bool) *Logger {
	return &Logger{
		minLevel:  level,
		isColored: useColors,
		logger:    log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelName := levelNames[level]
	
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	if l.isColored {
		color := levelColors[level]
		l.logger.Printf("%s %s%5s%s %s", 
			timestamp,
			color,
			levelName,
			colorReset,
			msg,
		)
	} else {
		l.logger.Printf("%s %5s %s",
			timestamp,
			levelName,
			msg,
		)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

// Logger interface methods
func Debug(format string, args ...interface{}) {
	defaultLogger.log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.log(INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.log(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.log(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.log(FATAL, format, args...)
} 