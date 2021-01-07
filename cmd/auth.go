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
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/meroxa/cli/credstore"

	"github.com/spf13/cobra"
)

const audience = "https://api.tjl.dev.meroxa.io/v1"
const clientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"
const domain = "auth.meroxa.io"
const meroxaLabel = "meroxa-cli"
const meroxaURL = "https://www.meroxa.io"
const scope = "user"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login or sign up to the Meroxa platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := login()
		if err != nil {
			return err
		}
		return nil
	},
}

// logoutCmd represents the logout command
//var logoutCmd = &cobra.Command{
//	Use:   "logout",
//	Short: "logout of the Meroxa platform",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		// TODO: add confirmation
//		err := logout()
//		if err != nil {
//			return err
//		}
//		fmt.Println("credentials cleared")
//		return nil
//	},
//}

func init() {
	// Login
	rootCmd.AddCommand(loginCmd)
	// Subcommands
	//// Logout
	//rootCmd.AddCommand(logoutCmd)
}

type deviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int64  `json:"interval"`
}

func getToken(domain string, deviceCode string, clientID string) (string, error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)

	fmt.Println("Printing device code in poling function :" + deviceCode)

	reqBody, err := json.Marshal(map[string]string{
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		"device_code": deviceCode,
		"client_id":   clientID,
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if len(body) > 0 {
		fmt.Println(string(body))
		var ff interface{}
		json.Unmarshal(body, &ff)
		result := ff.(map[string]interface{})
		accessToken := fmt.Sprintf("%v", result["access_token"])
		fmt.Println(accessToken)
		return accessToken, nil
	}

	return "", nil

}

func login() error {
	deviceUrl := fmt.Sprintf("https://%s/oauth/device/code", domain)
	reqBody, err := json.Marshal(map[string]string{
		"client_id": clientID,
		"scope":     scope,
		"audience":  audience,
	})
	if err != nil {
		log.Fatalln(err)
	}
	res, _ := http.Post(deviceUrl, "application/json", bytes.NewBuffer(reqBody))

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp deviceCodeResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("invalid device code response: %v", err)
		return err
	}
	verificationUrl := resp.VerificationURIComplete
	deviceCode := resp.DeviceCode

	log.Println(color.CyanString("You will now be taken to your browser for device authentication or open the url below in a browser."))
	log.Println(color.CyanString(verificationUrl))
	time.Sleep(1 * time.Second)
	err = openBrowser(verificationUrl)
	if err != nil {
		log.Println(color.RedString("Failed to open browser, you must open the link above to continue."))
	}
	var token string
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Pooling loop")
		token, err = getToken(domain, deviceCode, clientID)

		if err != nil {
			fmt.Print(err)
			continue
		}

		if len(token) > 5 {
			break
		}
	}

	err = credstore.Set(meroxaLabel, "meroxa-cli", token)
	if err != nil {
		fmt.Printf("Cant store token: %v", err)
		return err
	}

	return nil
}

func openBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}
