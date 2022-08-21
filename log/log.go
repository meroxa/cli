package log

import "os"

type Logger interface {
	LeveledLogger
	JSONLogger
	SpinnerLogger
}

type logger struct {
	LeveledLogger
	JSONLogger
	SpinnerLogger
}

func New(l1 LeveledLogger, l2 JSONLogger, l3 SpinnerLogger) Logger {
	return logger{
		LeveledLogger: l1,
		JSONLogger:    l2,
		SpinnerLogger: l3,
	}
}

func NewWithDevNull() Logger {
	o, _ := os.Open(os.DevNull)
	return New(NewLeveledLogger(o, Debug), NewJSONLogger(o), NewSpinnerLogger(o))
}
