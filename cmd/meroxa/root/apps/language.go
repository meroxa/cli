package apps

import (
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

const LanguageNotSupportedError = "Currently, we support \"javascript\", \"golang\", and \"python\""

// validateLanguage stops execution of the deployment in case language is not supported.
// It also consolidates lang used in API in case user specified a supported language using an unexpected description.
func validateLanguage(lang string) error {
	switch lang {
	case "go", turbine.GoLang:
		lang = turbine.GoLang
	case "js", turbine.JavaScript, turbine.NodeJs:
		lang = turbine.JavaScript
	case "py", turbine.Python3, turbine.Python:
		lang = turbine.Python
	case "rb", turbine.Ruby:
		lang = turbine.Ruby
	default:
		return fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
	}
	return nil
}
