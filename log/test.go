package log

import "bytes"

func NewTestLogger(level Level) *TestLogger {
	var buf bytes.Buffer
	return &TestLogger{
		buf: &buf,
		Logger: New(
			NewLeveledLogger(&buf, level),
			NewJSONLogger(&buf),
		),
	}
}

type TestLogger struct {
	Logger
	buf *bytes.Buffer
}

var _ Logger = (*TestLogger)(nil)

func (l *TestLogger) String() string {
	return l.buf.String()
}
