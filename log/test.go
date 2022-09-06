package log

import "bytes"

func NewTestLogger() *TestLogger {
	var leveledBuf bytes.Buffer
	var jsonBuf bytes.Buffer
	var spinnerBuf bytes.Buffer
	return &TestLogger{
		leveledBuf: &leveledBuf,
		jsonBuf:    &jsonBuf,
		spinnerBuf: &spinnerBuf,
		Logger: New(
			NewLeveledLogger(&leveledBuf, Debug),
			NewJSONLogger(&jsonBuf),
			NewSpinnerLogger(&spinnerBuf),
		),
	}
}

type TestLogger struct {
	Logger
	leveledBuf *bytes.Buffer
	jsonBuf    *bytes.Buffer
	spinnerBuf *bytes.Buffer
}

var _ Logger = (*TestLogger)(nil)

func (l *TestLogger) JSONOutput() string {
	return l.jsonBuf.String()
}

func (l *TestLogger) LeveledOutput() string {
	return l.leveledBuf.String()
}

func (l *TestLogger) SpinnerOutput() string {
	return l.spinnerBuf.String()
}
