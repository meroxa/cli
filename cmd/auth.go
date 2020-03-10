/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const MeroxaDirPath = ".config/meroxa"
const ConfigFileName = "meroxa.config"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "log into the Meroxa platform",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := saveCreds(args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Println("login saved")
		return nil
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout of the Meroxa platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: add confirmation
		err := clearCreds()
		if err != nil {
			return err
		}
		fmt.Println("credentials cleared")
		return nil
	},
}

var whoAmICmd = &cobra.Command{
	Use:   "whoami",
	Short: "retrieve currently logged in user",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _, err := readCreds()
		if err != nil {
			return err
		}
		fmt.Printf("username: %s", u)
		return nil
	},
}

func init() {
	// Login
	rootCmd.AddCommand(loginCmd)
	// Subcommands
	loginCmd.AddCommand(whoAmICmd)

	// Logout
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func saveCreds(username, password string) error {
	bytes := []byte(fmt.Sprintf("%s:%s", username, password))

	filePath, err := configFilePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readCreds() (string, string, error) {
	filePath, err := configFilePath()
	if err != nil {
		return "", "", err
	}
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}

	creds := strings.Split(string(dat), ":")
	return creds[0], creds[1], nil
}

func clearCreds() error {
	filePath, err := configFilePath()
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func createOrFindMeroxaConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	targetDirPath := home + "/" + MeroxaDirPath

	// Create Meroxa Config Dir if needed
	err = os.MkdirAll(targetDirPath, 0744)
	if err != nil {
		return "", err
	}

	return targetDirPath, nil

}

func configFilePath() (string, error) {
	mDir, err := createOrFindMeroxaConfigDir()
	if err != nil {
		return "", err
	}
	return mDir + "/" + ConfigFileName, nil
}
