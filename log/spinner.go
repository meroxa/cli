package log

import (
	"io"
	"log"
	"time"

	"github.com/briandowns/spinner"
)

type SpinnerLogger interface {
	StartSpinner(prefix, suffix string)
	StopSpinner(msg string)
}

func NewSpinnerLogger(out io.Writer) SpinnerLogger {
	return &spinnerLogger{l: log.New(out, "", 0)}
}

type spinnerLogger struct {
	l *log.Logger
	s *spinner.Spinner
}

func (l *spinnerLogger) StartSpinner(prefix, suffix string) {
	l.s = spinner.New(spinner.CharSets[14], 100*time.Millisecond) // nolint:gomnd
	l.s.Prefix = prefix
	l.s.Suffix = suffix
	l.s.Start()
}

func (l *spinnerLogger) StopSpinner(msg string) {
	l.s.Stop()
	l.l.Printf(msg)
}
