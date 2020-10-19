package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var outputFile *os.File

// StandardLogger constains underlying logger pkg.  Using this will allow for logging standardization later on
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger returns a logger pointer and the output file descriptor for safe clean up
func NewLogger() *StandardLogger {
	f, err1 := os.OpenFile("..\\src\\goapp1\\info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err1 != nil {
		fmt.Println("log failed")
		f.Close()
	}
	outputFile = f
	var newLogger = logrus.New()
	var StandardLogger = &StandardLogger{newLogger}
	StandardLogger.SetOutput(f)
	StandardLogger.SetFormatter(&logrus.JSONFormatter{}) // log can be set to json
	StandardLogger.Info("Logger Setup new")
	return StandardLogger

}

// CloseOutputFile is here to have this package close the output file
func CloseOutputFile() {
	/*
		d := syscall.Handle(outputFile.Fd()) // an attempt to solve blocking issue, leaving it here for posterity
		syscall.SetNonblock(d, true)
	*/
	err := outputFile.Close()
	if err != nil {
		fmt.Println("this file was already closed")
	} else {
		fmt.Println("output file is closed")
	}

}
