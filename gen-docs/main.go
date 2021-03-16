package main

import (
	"log"
	"os"

	cmd "github.com/meroxa/cli/cmd"

	"github.com/spf13/cobra/doc"
)

func main() {
	// set HOME env var so that default values involve user's home directory do not depend on the running user.
	os.Setenv("HOME", "/home/user")

	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs/cmd")
	if err != nil {
		log.Fatal(err)
	}
}
