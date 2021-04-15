package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/root"
	"github.com/spf13/cobra"
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

func renameFromUnderscoreToDash(dir string) {
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, err := file.Readdirnames(0) // 0 to read all files and folders
	if err != nil {
		log.Fatalf("failed reading directory: %s", err)
	}
	for _, name := range list {
		oldName := name
		newName := strings.Replace(oldName, "_", "-", -1)
		err := os.Rename(filepath.Join(dir, oldName), filepath.Join(dir, newName))
		if err != nil {
			log.Printf("error renaming file: %s", err)
			continue
		}
	}
}

func generateManPages(rootCmd *cobra.Command) error {
	header := &doc.GenManHeader{
		Title:   "Meroxa",
		Section: "1",
		// Source:  old2.VersionString(), TODO
		Manual: "Meroxa Manual",
	}

	return doc.GenManTree(rootCmd, header, "./etc/man/man1")
}

func generateMarkdownPages(rootCmd *cobra.Command) error {
	return doc.GenMarkdownTree(rootCmd, "./docs/cmd/md")
}

func generateDocsDotComPages(rootCmd *cobra.Command) error {
	filePrepender := func(filename string) string {
		createdAt := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		slug := strings.Replace(base, "_", "-", -1)
		url := "/cli/" + strings.ToLower(slug) + "/"
		return fmt.Sprintf(fmTemplate, createdAt, createdAt, strings.Replace(base, "_", " ", -1), slug, url)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		slug := strings.Replace(base, "_", "-", -1)
		return "/cli/" + strings.ToLower(slug) + "/"
	}

	return doc.GenMarkdownTreeCustom(rootCmd, "./docs/cmd/www", filePrepender, linkHandler)
}

func generateShellCompletionFiles(rootCmd *cobra.Command) {
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

func main() {
	// set HOME env var so that default values involve user's home directory do not depend on the running user.
	os.Setenv("HOME", "/home/user")

	rootCmd := root.Cmd()

	if err := generateManPages(rootCmd); err != nil {
		log.Fatal(err)
	}

	if err := generateMarkdownPages(rootCmd); err != nil {
		log.Fatal(err)
	}

	if err := generateDocsDotComPages(rootCmd); err != nil {
		log.Fatal(err)
	} else {
		renameFromUnderscoreToDash("./docs/cmd/www")
	}

	generateShellCompletionFiles(rootCmd)
}
