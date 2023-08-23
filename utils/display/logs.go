package display

import (
	"fmt"
	"time"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func LogsTable(ll *meroxa.Logs) string {
	var subTable string

	for i := len(ll.Data) - 1; i >= 0; i-- {
		l := ll.Data[i]
		subTable += fmt.Sprintf("[%s]\t%s\t%q\n", l.Timestamp.Format(time.RFC3339), l.Source, l.Log)
	}

	return subTable
}
