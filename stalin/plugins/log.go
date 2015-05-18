package plugins

import (
	"fmt"
	"runtime"
	. "stalin/log"
)

type Logger struct {
	name string
}

func NewLog(name string) *Logger {
	return &Logger{name: name}
}

func (l *Logger) Debug(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.Debug("[%s:%d(%s)] %s", file, line, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Info(format string, a ...interface{}) {
	Log.Info("[%s] %s", l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Warning(format string, a ...interface{}) {
	Log.Warning("[%s] %s", l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Error(format string, a ...interface{}) {
	Log.Error("[%s] %s", l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Fatal(format string, a ...interface{}) {
	Log.Fatal("[%s] %s\n", l.name, fmt.Sprintf(format, a...))
}
