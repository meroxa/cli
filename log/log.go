package log

type Logger interface {
	LeveledLogger
	JSONLogger
}

type logger struct {
	LeveledLogger
	JSONLogger
}

func New(l1 LeveledLogger, l2 JSONLogger) Logger {
	return logger{
		LeveledLogger: l1,
		JSONLogger:    l2,
	}
}
