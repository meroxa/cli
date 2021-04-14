package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	cmd "github.com/meroxa/cli/cmd"

	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
createdAt: %s
updatedAt: %s
title: "%s"
slug: %s
url: %s
---
`

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

	filePrepender := func(filename string) string {
		createdAt := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		url := "/cli/" + strings.ToLower(base) + "/"
		return fmt.Sprintf(fmTemplate, createdAt, createdAt, strings.Replace(base, "_", " ", -1), base, url)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return strings.ToLower(base)
	}

	// Generating Markdown Documents
	err := doc.GenMarkdownTreeCustom(rootCmd, "./docs/cmd", filePrepender, linkHandler)
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
