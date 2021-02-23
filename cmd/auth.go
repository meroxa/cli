/*
Copyright © 2020 Meroxa Inc

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
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"strings"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

const clientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"
const callbackURL = "http://localhost:21900/oauth/callback"
const domain = "auth.meroxa.io"
const audience = "https://api.tjl.dev.meroxa.io/v1"

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
	// Login
	RootCmd.AddCommand(loginCmd)
	// Logout
	RootCmd.AddCommand(logoutCmd)
}

func login() error {
	log.Println(color.CyanString("You will now be taken to your browser for authentication or open the url below in a browser."))
	authorizeUser(clientID, domain, callbackURL)
	return nil
}

// AuthorizeUser implements the PKCE OAuth2 flow.
func authorizeUser(clientID string, authDomain string, redirectURL string) {
	// initialize the code verifier
	var CodeVerifier, _ = cv.CreateCodeVerifier()

	// Create code_challenge with S256 method
	codeChallenge := CodeVerifier.CodeChallengeS256()

	// construct the authorization URL (with Auth0 as the authorization provider)
	authorizationURL := fmt.Sprintf(
		"https://%s/authorize?audience=%s"+
			"&scope=openid email offline_access user"+
			"&response_type=code&client_id=%s"+
			"&code_challenge=%s"+
			"&code_challenge_method=S256&redirect_uri=%s",
		authDomain, audience, clientID, codeChallenge, redirectURL)
	log.Println(color.CyanString(authorizationURL))

	// start a web server to listen on a callback URL
	server := &http.Server{Addr: redirectURL}

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		// get the authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Println("meroxa: Url Param 'code' is missing")
			io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		codeVerifier := CodeVerifier.String()
		accessToken, refreshToken, err := getAccessTokenAuth(clientID, codeVerifier, code, redirectURL)
		if err != nil {
			fmt.Println("meroxa: could not get access token")
			io.WriteString(w, "Error: could not retrieve tokens\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}
		cfg.Set("ACCESS_TOKEN", accessToken)
		cfg.Set("REFRESH_TOKEN", refreshToken)
		err = cfg.WriteConfig()
		if err != nil {
			fmt.Println("meroxa: could not write config file")
			io.WriteString(w, "Error: could not store access token\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// return an indication of success to the caller
		io.WriteString(w, `
		<html>
			<div style="height:100px; width:100%!; display:flex; flex-direction: column; justify-content: center; align-items:center; background-color:#2ecc71; color:white; font-size:22"><div>Success!</div></div>
		<p style="margin-top:20px; font-size:18; text-align:center">You are authenticated, you can now return to the program. This will auto-close</p>
		<script>window.onload=function(){setTimeout(this.close, 5000)}</script>
		</html>`)

		fmt.Println("Successfully logged in.")

		// close the HTTP server
		cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		fmt.Printf("meroxa: bad redirect URL: %s\n", err)
		os.Exit(1)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("meroxa: can't listen to port %s: %s\n", port, err)
		os.Exit(1)
	}

	// open a browser window to the authorizationURL
	err = open.Start(authorizationURL)
	if err != nil {
		fmt.Printf("meroxa: can't open browser to URL %s: %s\n", authorizationURL, err)
		os.Exit(1)
	}

	// start the blocking web server loop
	// this will exit when the handler gets fired and calls server.Close()
	server.Serve(l)
}

// getAccessToken trades the authorization code retrieved from the first OAuth2 leg for an access token
func getAccessTokenAuth(clientID string, codeVerifier string, authorizationCode string, callbackURL string) (string, string, error) {
	// set the url and form-encoded data for the POST to the access token endpoint
	url := "https://auth.meroxa.io/oauth/token"
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("meroxa: HTTP error: %s", err)
		return "", "", err
	}

	// process the response
	defer res.Body.Close()
	var responseData map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("meroxa: JSON error: %s", err)
		return "", "", err
	}

	// retrieve the access token out of the map, and return to caller
	accessToken := responseData["access_token"].(string)
	refreshToken := responseData["refresh_token"].(string)
	return accessToken, refreshToken, nil
}

// cleanup closes the HTTP server
func cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go server.Close()
}

func refresh(domain string, clientID string, refreshToken string) (accessToken string, newRefreshToken string, err error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)
	var tokenBody = make(map[string]string)
	tokenBody["client_id"] = clientID

	if refreshToken != "" {
		tokenBody["grant_type"] = "refresh_token"
		tokenBody["refresh_token"] = refreshToken
	} else {
		return "", "", errors.New("no refresh token")
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
		return "", fmt.Errorf("please login or signup by running 'meroxa login'")
	}

	// check token exp
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("please login or signup by running 'meroxa login'")
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
	accessToken, _, err = refresh(domain, clientID, refreshToken)
	if err != nil {
		return "", fmt.Errorf("please login or signup by running 'meroxa login'")
	}
	cfg.Set("ACCESS_TOKEN", accessToken)

	return accessToken, nil
}
