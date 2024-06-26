package log

import (
	"context"
	"encoding/json"
	"io"
	"log"
)

type JSONLogger interface {
	JSON(ctx context.Context, data interface{})
}

func NewJSONLogger(out io.Writer) JSONLogger {
	return &jsonLogger{l: log.New(out, "", 0)}
}

type jsonLogger struct {
	l *log.Logger
}

func (l *jsonLogger) JSON(_ context.Context, data interface{}) {
	if raw, ok := data.(string); ok {
		l.l.Print(raw)
		return
	}

	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		l.l.Printf("could not marshal JSON: %s", err.Error())
		return
	}
	l.l.Print(string(p))
}
