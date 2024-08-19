package log

import (
	"context"
	"github.com/tryfix/log"
)

type Logger struct {
	log log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		log: log.Constructor.Log(
			log.WithColors(true),
			log.WithLevel(log.TRACE),
			log.WithFilePath(true),
			log.WithSkipFrameCount(3),
		),
	}
}

func (l *Logger) Fatal(message interface{}, params ...interface{}) {
	l.log.Fatal(message, params)
}

func (l *Logger) Error(message interface{}, params ...interface{}) {
	l.log.Error(message, params)
}

func (l *Logger) Warn(message interface{}, params ...interface{}) {
	l.log.Warn(message, params)
}

func (l *Logger) Debug(message interface{}, params ...interface{}) {
	l.log.Debug(message, params)
}

func (l *Logger) Info(message interface{}, params ...interface{}) {
	l.log.Info(message, params)
}

func (l *Logger) Trace(message interface{}, params ...interface{}) {
	l.log.Trace(message, params)
}

func (l *Logger) FatalContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.log.FatalContext(ctx, message, params)
}

func (l *Logger) ErrorContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.log.ErrorContext(ctx, message, params)
}

func (l *Logger) WarnContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.log.WarnContext(ctx, message, params)
}

func (l *Logger) DebugContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.log.DebugContext(ctx, message, params)
}

func (l *Logger) InfoContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.log.InfoContext(ctx, message, params)
}

func (l *Logger) TraceContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.log.TraceContext(ctx, message, params)
}
