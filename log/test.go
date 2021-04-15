package log

import "bytes"

func NewTestLogger() *TestLogger {
	var leveledBuf bytes.Buffer
	var jsonBuf bytes.Buffer
	return &TestLogger{
		leveledBuf: &leveledBuf,
		jsonBuf:    &jsonBuf,
		Logger: New(
			NewLeveledLogger(&leveledBuf, Debug),
			NewJSONLogger(&jsonBuf),
		),
	}
}

type TestLogger struct {
	Logger
	leveledBuf *bytes.Buffer
	jsonBuf    *bytes.Buffer
}

var _ Logger = (*TestLogger)(nil)

func (l *TestLogger) JSONOutput() string {
	return l.jsonBuf.String()
}

func (l *TestLogger) LeveledOutput() string {
	return l.leveledBuf.String()
}
