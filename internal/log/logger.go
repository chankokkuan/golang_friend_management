package log

import "context"

type Logger interface {
	Attach(c context.Context, keyvals ...interface{}) context.Context
	With(keyvals ...interface{}) Logger
	WithCtx(c context.Context) Logger
	Debug(keyvals ...interface{})
	Info(keyvals ...interface{})
	Warn(keyvals ...interface{})
	Error(keyvals ...interface{})
}
