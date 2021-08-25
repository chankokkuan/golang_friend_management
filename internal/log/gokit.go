package log

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

type contextKeyType struct{}

var contextKey contextKeyType

type gokitWrapper struct {
	logger log.Logger
}

func NewGokitWrapper(serviceName string) Logger {
	return &gokitWrapper{newLogger(serviceName)}
}

func newLogger(serviceName string) log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger,
		"service", serviceName,
		"time", log.DefaultTimestampUTC,
		"caller", log.Caller(4),
	)
	return logger
}

func (w *gokitWrapper) With(keyvals ...interface{}) Logger {
	logger := log.With(w.logger, keyvals...)
	return &gokitWrapper{logger}
}

func (w *gokitWrapper) Attach(c context.Context, keyvals ...interface{}) context.Context {
	wrapper := w.WithCtx(c)
	wrapper = wrapper.With(keyvals...)
	return context.WithValue(c, contextKey, wrapper)
}

func (w *gokitWrapper) WithCtx(c context.Context) Logger {
	if logger, ok := c.Value(contextKey).(*gokitWrapper); ok && logger != nil {
		return logger
	}
	return w
}

func (w *gokitWrapper) Debug(keyvals ...interface{}) {
	level.Debug(w.logger).Log(keyvals...)
}

func (w *gokitWrapper) Info(keyvals ...interface{}) {
	level.Info(w.logger).Log(keyvals...)
}

func (w *gokitWrapper) Warn(keyvals ...interface{}) {
	level.Warn(w.logger).Log(keyvals...)
}

func (w *gokitWrapper) Error(keyvals ...interface{}) {
	level.Error(w.logger).Log(keyvals...)
}
