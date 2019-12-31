package logger

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger allows to log.
type Logger struct {
	logrus          *logrus.Logger
	mu              *sync.RWMutex
	sensitiveFields []string
}

// Options overrides brhavior of Logger
type Option func(*Logger) error

// WithSensitiveFields allows to set sensitive fields
// which will be filtered.
func WithSensitiveFields(fields []string) Option {
	f := func(l *Logger) error {
		l.sensitiveFields = fields
		return nil
	}
	return f
}

// SetFormatter set logrus formatter
func SetFormatter(f logrus.Formatter) Option {
	return func(l *Logger) error {
		l.logrus.Formatter = f
		return nil
	}
}

// New initialize logger with dependencies.
func New(options ...Option) (*Logger, error) {
	l := Logger{
		logrus:          logrus.New(),
		mu:              new(sync.RWMutex),
		sensitiveFields: make([]string, 0),
	}

	l.logrus.Formatter = &logrus.JSONFormatter{}

	for _, option := range options {
		if err := option(&l); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	return &l, nil
}

// Info prints at info level.
func (l *Logger) Info(info string, extra map[string]interface{}) {
	l.mu.RLock()
	{
		l.filterFields(extra)
		l.logrus.WithFields(logrus.Fields(extra)).
			Info(info)
	}
	l.mu.RUnlock()
}

// Warn prints at warn level.
func (l *Logger) Warn(warn string, extra map[string]interface{}) {
	l.mu.RLock()
	{
		l.filterFields(extra)
		l.logrus.WithFields(logrus.Fields(extra)).
			Warn(warn)
	}
	l.mu.RUnlock()
}

// Error prints at error level.
func (l *Logger) Error(err error, extra map[string]interface{}) {
	l.mu.RLock()
	{
		l.filterFields(extra)
		l.logrus.WithFields(logrus.Fields(extra)).
			Error(err.Error())
	}
	l.mu.RUnlock()
}

// Fatal prints at fatal level.
func (l *Logger) Fatal(err error, extra map[string]interface{}) {
	l.mu.RLock()
	{
		l.filterFields(extra)
		l.logrus.WithFields(logrus.Fields(extra)).
			Fatal(err)
	}
	l.mu.RUnlock()
}

func (l *Logger) filterFields(fields map[string]interface{}) {
	for _, field := range l.sensitiveFields {
		if _, ok := fields[field]; ok {
			fields[field] = "[FILTERED]"
		}
	}
}
