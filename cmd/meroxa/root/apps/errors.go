package apps

import (
	"fmt"

	"github.com/meroxa/turbine-core/pkg/ir"
)

const (
	resourceInvalidError = `⚠️  Run 'meroxa resources list' to verify that the resource names ` +
		`defined in your Turbine app are identical to the resources you have ` +
		`created on the Meroxa Platform before deploying again`
)

func newLangUnsupportedError(lang ir.Lang) error {
	return fmt.Errorf(
		`language %q not supported. `+
			`supported languages "javascript", "golang", "python", "ruby", and "java"`,
		lang,
	)
}

func wrapErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	err := errs[0]
	for _, e := range errs[1:] {
		err = fmt.Errorf("%w; %v", err, e)
	}

	return err
}
