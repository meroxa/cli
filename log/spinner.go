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
	Warning    = "warning"
)

type SpinnerLogger interface {
	StartSpinner(prefix, suffix string)
	StopSpinner(msg string)
	StopSpinnerWithStatus(msg, status string)
	SuccessfulCheck() string
	FailedMark() string
}

func NewSpinnerLogger(out io.Writer) SpinnerLogger {
	return &spinnerLogger{l: log.New(out, "", 0), out: out}
}

type spinnerLogger struct {
	l   *log.Logger
	s   *spinner.Spinner
	out io.Writer
}

func (l *spinnerLogger) StartSpinner(prefix, suffix string) {
	l.s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(l.out))
	l.s.Prefix = prefix
	l.s.Suffix = " " + suffix
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
	} else if status == Warning {
		msg = fmt.Sprintf("\t%s %s", l.WarningMark(), msg)
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

func (l *spinnerLogger) WarningMark() string {
	return color.New(color.FgYellow).Sprintf("⚡")
}
