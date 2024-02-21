package apps

import (
	"fmt"

	"github.com/meroxa/turbine-core/v2/pkg/ir"
)

func newLangUnsupportedError(lang ir.Lang) error {
	return fmt.Errorf(
		`language %q not supported. `+
			`supported languages "javascript", "golang", "python", and "ruby (beta)"`,
		lang,
	)
}
