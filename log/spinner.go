package log

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/fatih/color"

	"github.com/briandowns/spinner"
)

const (
	Successful = "successful"
	Failed     = "failed"
)

type SpinnerLogger interface {
	StartSpinner(prefix, suffix string)
	StopSpinner(msg string)
	StopSpinnerWithStatus(msg, status string)
	SuccessfulCheck() string
	FailedMark() string
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

func (l *spinnerLogger) StopSpinnerWithStatus(msg, status string) {
	if status == Failed {
		msg = fmt.Sprintf("\t%s %s", l.FailedMark(), msg)
	} else if status == Successful {
		msg = fmt.Sprintf("\t%s %s", l.SuccessfulCheck(), msg)
	}
	l.s.Stop()
	l.l.Printf(msg)
}

func (l *spinnerLogger) SuccessfulCheck() string {
	return color.New(color.FgGreen).Sprintf("✔")
}

func (l *spinnerLogger) FailedMark() string {
	return color.New(color.FgRed).Sprintf("x")
}
