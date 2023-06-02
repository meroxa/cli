package apps

import (
	"fmt"

	"github.com/meroxa/turbine-core/pkg/ir"
)

func newLangUnsupportedError(lang ir.Lang) error {
	return fmt.Errorf(
		`language %q not supported. `+
			`supported languages "javascript", "golang", "python", "ruby", and "java"`,
		lang,
	)
}
