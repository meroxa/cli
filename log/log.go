package log

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
