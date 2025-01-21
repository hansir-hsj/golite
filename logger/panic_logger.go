package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type PanicLogger struct {
	filePath string
	file     *os.File
}

func NewPanicLogger(ctx context.Context, confDir ...string) (*PanicLogger, error) {
	var filePath string

	if len(confDir) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		filePath = filepath.Join(dir, "/log/panic.log")
	} else {
		loggerConfig := filepath.Join(confDir[0], LoggerConfigFile)
		logConf, err := parse(loggerConfig)
		if err != nil {
			return nil, err
		}
		filePath = filepath.Join(logConf.Dir, "panic.log")
	}

	target, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &PanicLogger{
		filePath: filePath,
		file:     target,
	}, nil
}

func (l *PanicLogger) caller() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return ""
	}
	return strings.Join([]string{file, strconv.Itoa(line)}, ":")
}

func (l *PanicLogger) Report(ctx context.Context, p any) {
	msg := fmt.Sprintf("Recover from panic: %v", p)
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	stack = stack[:length]

	fmt.Fprintf(l.file, "%s\n%s\nStack:\n%s\n", msg, l.caller(), stack)
}