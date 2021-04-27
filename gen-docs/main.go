package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

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

func renameFromUnderscoreToDash(dir string) error {
	file, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("failed opening directory: %w", err)
	}
	defer file.Close()

	list, err := file.Readdirnames(0) // 0 to read all files and folders
	if err != nil {
		return fmt.Errorf("failed reading directory: %w", err)
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
	return nil
}

func generateManPages(rootCmd *cobra.Command) error {
	header := &doc.GenManHeader{
		Title:   "Meroxa",
		Section: "1",
		Source:  "Meroxa CLI ", // TODO change this to actually get version
		Manual:  "Meroxa Manual",
	}

	return doc.GenManTree(rootCmd, header, "./etc/man/man1")
}

func generateMarkdownPages(rootCmd *cobra.Command) error {
	return doc.GenMarkdownTree(rootCmd, "./docs/cmd/md")
}

func generateDocsDotComPages(rootCmd *cobra.Command) error {
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		slug := strings.Replace(base, "_", "-", -1)
		url := "/cli/" + strings.ToLower(slug) + "/"
		return fmt.Sprintf(fmTemplate, "", "", strings.Replace(base, "_", " ", -1), slug, url)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		slug := strings.Replace(base, "_", "-", -1)
		return "/cli/" + strings.ToLower(slug) + "/"
	}

	return doc.GenMarkdownTreeCustom(rootCmd, "./docs/cmd/www", filePrepender, linkHandler)
}

func generateShellCompletionFiles(rootCmd *cobra.Command) error {
	if err := rootCmd.GenBashCompletionFile("./etc/completion/meroxa.completion.sh"); err != nil {
		return fmt.Errorf("could not generate bash completion: %w", err)
	}
	if err := rootCmd.GenZshCompletionFile("./etc/completion/meroxa.completion.zsh"); err != nil {
		return fmt.Errorf("could not generate zsh completion: %w", err)
	}
	if err := rootCmd.GenFishCompletionFile("./etc/completion/meroxa.completion.fish", true); err != nil {
		return fmt.Errorf("could not generate fish completion: %w", err)
	}
	if err := rootCmd.GenPowerShellCompletionFile("./etc/completion/meroxa.completion.ps1"); err != nil {
		return fmt.Errorf("could not generate powershell completion: %w", err)
	}
	return nil
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
	}

	if err := renameFromUnderscoreToDash("./docs/cmd/www"); err != nil {
		log.Fatal(err)
	}

	if err := generateShellCompletionFiles(rootCmd); err != nil {
		log.Fatal(err)
	}
}
