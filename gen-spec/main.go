package main

import (
	"os"

	"github.com/meroxa/cli/cmd/meroxa/root"

	fig "github.com/withfig/autocomplete-tools/integrations/cobra"
)

func main() {
	// set HOME env var so that default values involving the user's home directory do not depend on the running user.
	os.Setenv("HOME", "/home/user")

	rootCmd := root.Cmd()
	spec := fig.GenerateCompletionSpec(rootCmd)
	f, err := os.Create("./spec/meroxa.ts")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(spec.ToTypescript())
	if err != nil {
		panic(err)
	}
}
