package rlog

import (
	"fmt"
	"os"
	"time"

	"github.com/raxisau/ds"
)

// This is the debug levels for the output
const (
	ERRORLEVEL = iota
	WARNINGLEVEL
	NOTICELEVEL
	INFOLEVEL
	DEBUGLEVEL
	TRACELEVEL

	TRACELEVELNAME    = "TRACE"
	DEBUGLEVELNAME    = "DEBUG"
	INFOLEVELNAME     = "INFO"
	NOTICELEVELNAME   = "NOTICE"
	WARNINGLEVELNAME  = "WARNING"
	ERRORLEVELNAME    = "ERROR"
	CRITICALLEVELNAME = "CRITICAL"

	logScrollBufferSize  = 1000
	logChannelBufferSize = 1000
)

type logRecord struct {
	logTime  time.Time
	logArgs  []interface{}
	logLevel string
}

var (
	logEventChannel chan *logRecord
	logOutFile      *os.File
	logClosed       bool = false
	loggingLevel    int  = INFOLEVEL
	scrollBuff           = ds.NewScrollBuffer(500)
)

// LogSetup Name says it all
func LogSetup(debugFile, logLevel string) {
	var err error

	if debugFile == "stdout" {
		logOutFile = os.Stdout
	} else if debugFile == "stderr" {
		logOutFile = os.Stderr
	} else {
		logOutFile, err = os.OpenFile(debugFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logOutFile = os.Stdout
		}
	}

	logEventChannel = make(chan *logRecord, logChannelBufferSize)
	go consumeLogEvents()
	SetLogLevel(logLevel)
}

func logSend(level string, args ...interface{}) {
	if logClosed {
		return
	}

	select {
	case logEventChannel <- &logRecord{logLevel: level, logTime: time.Now().UTC(), logArgs: args}:
	default:
	}
}

// Close the Log Output channel and close the file
func Close() {
	logClosed = true
	logEventChannel <- nil
}

func consumeLogEvents() {
	for {
		logRec := <-logEventChannel
		if logRec == nil {
			break
		}

		logLine := logRec.logTime.Format("15:04:05.000") + " " + logRec.logLevel
		for _, arg := range logRec.logArgs {
			logLine += " " + fmt.Sprint(arg)
		}
		scrollBuff.Put(logLine)
		logOutFile.WriteString(logLine)
		logOutFile.WriteString("\n")
	}

	if logOutFile != os.Stdout && logOutFile != os.Stderr {
		logOutFile.Close()
	}
	close(logEventChannel)
}

// GetLogLevel Name says it all
func GetLogLevel() string {
	switch loggingLevel {
	case TRACELEVEL:
		return TRACELEVELNAME
	case DEBUGLEVEL:
		return DEBUGLEVELNAME
	case INFOLEVEL:
		return INFOLEVELNAME
	case NOTICELEVEL:
		return NOTICELEVELNAME
	case ERRORLEVEL:
		return ERRORLEVELNAME
	case WARNINGLEVEL:
		fallthrough
	default:
		return WARNINGLEVELNAME
	}
}

// SetLogLevel Name says it all
func SetLogLevel(logLevel string) string {
	switch logLevel {
	case TRACELEVELNAME:
		loggingLevel = TRACELEVEL
		logSend(INFOLEVELNAME, "Logging Level set to:", logLevel)

	case DEBUGLEVELNAME:
		loggingLevel = DEBUGLEVEL
		logSend(INFOLEVELNAME, "Logging Level set to:", logLevel)

	case INFOLEVELNAME:
		loggingLevel = INFOLEVEL
		logSend(INFOLEVELNAME, "Logging Level set to:", logLevel)

	case NOTICELEVELNAME:
		loggingLevel = NOTICELEVEL
		logSend(INFOLEVELNAME, "Logging Level set to:", logLevel)

	case ERRORLEVELNAME:
		loggingLevel = ERRORLEVEL
		logSend(INFOLEVELNAME, "Logging Level set to:", logLevel)

	case WARNINGLEVELNAME:
		fallthrough
	default:
		loggingLevel = WARNINGLEVEL
		logLevel = WARNINGLEVELNAME
		logSend(INFOLEVELNAME, "Logging Level set to:", WARNINGLEVELNAME)
	}
	return logLevel
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	logSend(CRITICALLEVELNAME, args...)
	Close()
	time.Sleep(time.Second)
	os.Exit(1)
}

// Critical logs a message using CRITICAL as log level.
func Critical(args ...interface{}) {
	logSend(CRITICALLEVELNAME, args...)
}

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
	if loggingLevel >= ERRORLEVEL {
		logSend(ERRORLEVELNAME, args...)
	}
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
	if loggingLevel >= WARNINGLEVEL {
		logSend(WARNINGLEVELNAME, args...)
	}
}

// Notice logs a message using NOTICE as log level.
func Notice(args ...interface{}) {
	if loggingLevel >= NOTICELEVEL {
		logSend(NOTICELEVELNAME, args...)
	}
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
	if loggingLevel >= INFOLEVEL {
		logSend(INFOLEVELNAME, args...)
	}
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
	if loggingLevel >= DEBUGLEVEL {
		logSend(DEBUGLEVELNAME, args...)
	}
}

// Trace logs a message using TRACE as log level.
func Trace(args ...interface{}) {
	if loggingLevel >= TRACELEVEL {
		logSend(TRACELEVELNAME, args...)
	}
}

// GetTail gets the last 500 lines of the debugging output
func GetTail() []string {
	return scrollBuff.GetBuffer()
}
