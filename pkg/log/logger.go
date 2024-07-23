package log

import "github.com/tryfix/log"

type Logger struct {
	log log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		log: log.Constructor.Log(
			log.WithColors(true),
			log.WithLevel("DEBUG"),
			log.WithFilePath(true),
			log.WithSkipFrameCount(4),
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
