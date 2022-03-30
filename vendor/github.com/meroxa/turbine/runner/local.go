//go:build !platform
// +build !platform

package runner

import (
	"log"

	"github.com/meroxa/turbine"

	"github.com/meroxa/turbine/local"
)

func Start(app turbine.App) {
	lv := local.New()
	err := app.Run(lv)
	if err != nil {
		log.Fatalln(err)
	}
}
