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
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const MeroxaDirPath = ".config/meroxa"
const ConfigFileName = "meroxa.config"

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "sign up to the Meroxa platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := prompt("Username", usernameValidator, false)
		if err != nil {
			fmt.Println("Username invalid: ", err)
			return err
		}
		p, err := prompt("Password", passwordValidator, true)
		if err != nil {
			fmt.Println("Password invalid: ", err)
			return err
		}
		e, err := prompt("Email", emailValidator, false)
		if err != nil {
			fmt.Println("Email invalid: ", err)
			return err
		}

		fmt.Println("Registering user...")
		err = signup(u, p, e)
		if err != nil {
			return err
		}

		fmt.Println("User registered!")
		return nil
	},
}

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

	// Signup
	rootCmd.AddCommand(signupCmd)

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

// Prompts

func passwordValidator(input string) error {
	if len(input) < 8 || len(input) > 256 {
		return errors.New("password should be between 8 and 256 characters long")
	}
	return nil
}

func emailValidator(input string) error {
	rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(input) > 254 || !rxEmail.MatchString(input) {
		return errors.New("email provided is invalid")
	}

	return nil
}

func usernameValidator(input string) error {
	usernameRegexp := regexp.MustCompile(`(?i)^[a-z][a-z0-9]{2,11}$`)
	if len(input) < 3 {
		return errors.New("input should be at least 3 characters long")
	}

	if len(input) > 12 {
		return errors.New("input should be no longer than 12 characters")
	}

	if !usernameRegexp.Match([]byte(input)) {
		return errors.New("username should start only contain alphanumeric characters and start with a letter")
	}

	return nil
}

func prompt(label string, validator func(input string) error, mask bool) (string, error) {
	p := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}

	if mask {
		p.Mask = '*'
	}

	return p.Run()
}

func signup(username, password, email string) error {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	b64Password := b64.StdEncoding.EncodeToString([]byte(password))

	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": b64Password,
		"email":    email,
	})
	if err != nil {
		return err
	}

	apiURL := "https://api.meroxa.io/v1/"
	if u := os.Getenv("API_URL"); u != "" {
		apiURL = u
	}
	resp, err := c.Post(apiURL+"users", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return fmt.Errorf("error %+v", string(body))
	}

	return nil
}
