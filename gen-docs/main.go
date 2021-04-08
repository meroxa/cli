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

	rootCmd := cmd.RootCmd()

	// Generating Man Pages
	header := &doc.GenManHeader{
		Title:   "Meroxa",
		Section: "1",
		Source:  cmd.VersionString(),
		Manual:  "Meroxa Manual",
	}

	if err := doc.GenManTree(rootCmd, header, "./etc/man/man1"); err != nil {
		log.Fatal(err)
	}

	// Generating Markdown Documents
	err := doc.GenMarkdownTree(rootCmd, "./docs/cmd")
	if err != nil {
		log.Fatal(err)
	}

	// Generating Bash Completion File
	if err := rootCmd.GenBashCompletionFile("./etc/completion/meroxa.completion.sh"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.GenZshCompletionFile("./etc/completion/meroxa.completion.zsh"); err != nil {
		log.Fatal(err)
	}

	// Generating Fish Completion File
	fishCompletionFile, err := os.Create("./etc/completion/meroxa.completion.fish")

	if err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.GenFishCompletion(fishCompletionFile, true); err != nil {
		log.Fatal(err)
	}

	// Generating Power Shell File
	powerShellCompletionFile, err := os.Create("./etc/completion/meroxa.completion.ps1")

	if err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.GenPowerShellCompletion(powerShellCompletionFile); err != nil {
		log.Fatal(err)
	}
}
