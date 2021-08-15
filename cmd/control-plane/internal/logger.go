package internal

import (
	_log "log"
	"os"
)

type logger struct {
	l *_log.Logger
}

func NewLogger() *logger {
	return &logger{
		l: _log.New(os.Stdout, "[cache] ", _log.LstdFlags),
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.l.Printf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}
