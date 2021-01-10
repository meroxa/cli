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
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

const audience = "https://api.tjl.dev.meroxa.io/v1"
const clientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"
const domain = "auth.meroxa.io"
const scope = "user offline_access"

// loginCmd represents the login command
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
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout of the Meroxa platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: add confirmation
		cfg.Set("ACCESS_TOKEN", "")
		cfg.Set("REFRESH_TOKEN", "")
		err := cfg.WriteConfig()
		if err != nil {
			return err
		}
		fmt.Println("credentials cleared")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("keychain", "", "", "keychain name, uses login by default")
	rootCmd.AddCommand(logoutCmd)
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

	var (
		accessToken  string
		refreshToken string
	)

	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Pooling auth endpoint")
		accessToken, refreshToken, err = getTokens(domain, deviceCode, clientID, "")

		if err != nil {
			fmt.Print(err)
			continue
		}

		if len(accessToken) > 5 {
			break
		}
	}

	cfg.Set("ACCESS_TOKEN", accessToken)
	cfg.Set("REFRESH_TOKEN", refreshToken)
	err = cfg.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing config: %v", err)
		return err
	}

	return nil
}

func getTokens(domain string, deviceCode string, clientID string, refreshToken string) (accessToken string, newRefreshToken string, err error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)
	var tokenBody = make(map[string]string)
	tokenBody["client_id"] = clientID

	if refreshToken != "" {
		tokenBody["grant_type"] = "refresh_token"
		tokenBody["refresh_token"] = refreshToken
	} else {
		tokenBody["grant_type"] = "urn:ietf:params:oauth:grant-type:device_code"
		tokenBody["device_code"] = deviceCode
	}
	requestBody, err := json.Marshal(tokenBody)
	if err != nil {
		fmt.Println("marshal:", err)
		return "", "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("request error:", err)
		return "", "", err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if len(responseBody) > 0 {
		var ff interface{}
		json.Unmarshal(responseBody, &ff)
		result := ff.(map[string]interface{})
		accessToken = fmt.Sprintf("%v", result["access_token"])
		if refreshToken == "" {
			newRefreshToken = fmt.Sprintf("%v", result["refresh_token"])
		}

		return accessToken, newRefreshToken, nil
	}

	return "", "", nil
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

func getAccessToken() (string, error) {
	// check access token expiration
	accessToken := cfg.GetString("ACCESS_TOKEN")
	if accessToken == "" {
		return "", fmt.Errorf("please login or signup by running 'meroxa login'")
	}

	// check access exp and grab refresh
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("2please login or signup by running 'meroxa login'")
	}

	// check token exp
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("3please login or signup by running 'meroxa login'")
	}

	var exp time.Time
	tokenExp := claims["exp"].(float64)
	exp = time.Unix(int64(tokenExp), 0)

	if exp.After(time.Now()) {
		return accessToken, nil
	}

	// access token is expire, use refresh
	refreshToken := cfg.GetString("REFRESH_TOKEN")
	if refreshToken == "" {
		return "", fmt.Errorf("please login or signup by running 'meroxa login'")
	}
	accessToken, _, err = getTokens(domain, "", clientID, refreshToken)
	if err != nil {
		return "", fmt.Errorf("please login or signup by running 'meroxa login'")
	}

	return accessToken, nil
}

type deviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int64  `json:"interval"`
}
