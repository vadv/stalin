package log

import (
	"fmt"
)

var Log = new(logger)

type loggerInterface interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type logger struct {
	Logger loggerInterface
}

func (l *logger) SetLogger(i loggerInterface) {
	l.Logger = i
}

func (l *logger) Debug(format string, a ...interface{}) {
	l.Logger.Debugf("%s", fmt.Sprintf(format, a...))
}

func (l *logger) Info(format string, a ...interface{}) {
	l.Logger.Infof("%s", fmt.Sprintf(format, a...))
}

func (l *logger) Warning(format string, a ...interface{}) {
	l.Logger.Warnf("%s", fmt.Sprintf(format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	l.Logger.Errorf("%s", fmt.Sprintf(format, a...))
}

func (l *logger) Fatal(format string, a ...interface{}) {
	l.Logger.Fatalf("%s", fmt.Sprintf(format, a...))
}
