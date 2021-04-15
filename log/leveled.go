package log

import (
	"context"
	"io"
	"log"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type LeveledLogger interface {
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, fmt string, args ...interface{})
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, fmt string, args ...interface{})
	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, fmt string, args ...interface{})
	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, fmt string, args ...interface{})
}

func NewLeveledLogger(out io.Writer, level Level) LeveledLogger {
	return &leveledLogger{
		l:     log.New(out, "", 0),
		level: level,
	}
}

type leveledLogger struct {
	l     *log.Logger
	level Level
}

func (l *leveledLogger) Debug(ctx context.Context, msg string) {
	if !l.checkLevel(Debug) {
		return // skip
	}
	l.l.Print(msg)
}

func (l *leveledLogger) Debugf(ctx context.Context, fmt string, args ...interface{}) {
	if !l.checkLevel(Debug) {
		return // skip
	}
	l.l.Printf(fmt, args...)
}

func (l *leveledLogger) Info(ctx context.Context, msg string) {
	if !l.checkLevel(Info) {
		return // skip
	}
	l.l.Print(msg)
}

func (l *leveledLogger) Infof(ctx context.Context, fmt string, args ...interface{}) {
	if !l.checkLevel(Info) {
		return // skip
	}
	l.l.Printf(fmt, args...)
}

func (l *leveledLogger) Warn(ctx context.Context, msg string) {
	if !l.checkLevel(Warn) {
		return // skip
	}
	l.l.Print(msg)
}

func (l *leveledLogger) Warnf(ctx context.Context, fmt string, args ...interface{}) {
	if !l.checkLevel(Warn) {
		return // skip
	}
	l.l.Printf(fmt, args...)
}

func (l *leveledLogger) Error(ctx context.Context, msg string) {
	if !l.checkLevel(Error) {
		return // skip
	}
	l.l.Print(msg)
}

func (l *leveledLogger) Errorf(ctx context.Context, fmt string, args ...interface{}) {
	if !l.checkLevel(Error) {
		return // skip
	}
	l.l.Printf(fmt, args...)
}

func (l *leveledLogger) checkLevel(level Level) bool {
	return l.level <= level
}
