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
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/meroxa/cli/credstore"

	rndm "github.com/nmrshll/rndm-go"
	"github.com/spf13/cobra"
)

const ClientID = "780birp93au255rqsenjc2o0pu"
const CallbackURL = "http://localhost:21900/oauth/callback"
const AuthDomain = "tjl-meroxa.auth.us-east-1.amazoncognito.com"
const meroxaLabel = "meroxa-cli"
const meroxaURL = "https://www.meroxa.io"
const PORT = 21900
const authTimeout = 300
const oauthStateStringContextKey = 232

var oauthConfig = &oauth2.Config{
	ClientID: ClientID, // also known as client key sometimes
	Scopes:   []string{"email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://tjl-meroxa.auth.us-east-1.amazoncognito.com/oauth2/authorize",
		TokenURL: "https://tjl-meroxa.auth.us-east-1.amazoncognito.com/oauth2/token",
	},
	RedirectURL: CallbackURL,
}

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

//var whoAmICmd = &cobra.Command{
//	Use:   "whoami",
//	Short: "retrieve currently logged in user",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		u, _, err := readCreds()
//		if err != nil {
//			return err
//		}
//		fmt.Printf("username: %s", u)
//		return nil
//	},
//}

func init() {
	// Login
	rootCmd.AddCommand(loginCmd)
	// Subcommands
	//loginCmd.AddCommand(whoAmICmd)
	//// Logout
	//rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func login() error {

	ctx := context.Background()

	// Some random string, random for each request
	oauthStateString := rndm.String(8)
	ctx = context.WithValue(ctx, oauthStateStringContextKey, oauthStateString)
	urlString := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

	labelChan, stopHTTPServerChan, cancelAuthentication := startHTTPServer(ctx, oauthConfig)
	log.Println(color.CyanString("You will now be taken to your browser for authentication or open the url below in a browser."))
	log.Println(color.CyanString(urlString))
	log.Println(color.CyanString("If you are opening the url manually on a different machine you will need to curl the result url on this machine manually."))
	time.Sleep(1000 * time.Millisecond)
	err := openBrowser(urlString)
	if err != nil {
		log.Println(color.RedString("Failed to open browser, you MUST do the manual process."))
	}
	time.Sleep(600 * time.Millisecond)

	// shutdown the server after timeout
	go func() {
		log.Printf("Authentication will be cancelled in %s seconds", strconv.Itoa(authTimeout))
		time.Sleep(authTimeout * time.Second)
		stopHTTPServerChan <- struct{}{}
	}()

	select {
	// wait for client on clientChan
	case label := <-labelChan:
		// After the callbackHandler returns a client, it's time to shutdown the server gracefully
		stopHTTPServerChan <- struct{}{}
		fmt.Println("debug", *label)
		user, token, err := credstore.Get(meroxaLabel, meroxaURL)
		if err != nil {
			return err
		}
		fmt.Printf("User: %s Token: %s", user, token)

		return nil

		// if authentication process is cancelled first return an error
	case <-cancelAuthentication:
		return fmt.Errorf("authentication timed out and was cancelled")
	}

	return nil
}

func startHTTPServer(ctx context.Context, conf *oauth2.Config) (tokenChan chan *string, stopHTTPServerChan chan struct{}, cancelAuthentication chan struct{}) {
	// init returns
	tokenChan = make(chan *string)
	stopHTTPServerChan = make(chan struct{})
	cancelAuthentication = make(chan struct{})

	http.HandleFunc("/oauth/callback", callbackHandler(ctx, conf, tokenChan))
	srv := &http.Server{Addr: ":" + strconv.Itoa(PORT)}

	// handle server shutdown signal
	go func() {
		// wait for signal on stopHTTPServerChan
		<-stopHTTPServerChan
		log.Println("Shutting down server...")

		// give it 5 sec to shutdown gracefully, else quit program
		d := time.Now().Add(5 * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf(color.RedString("Auth server could not shutdown gracefully: %v"), err)
		}

		// after server is shutdown, quit program
		cancelAuthentication <- struct{}{}
	}()

	// handle callback request
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		fmt.Println("Server gracefully stopped")
	}()

	return tokenChan, stopHTTPServerChan, cancelAuthentication
}

func callbackHandler(ctx context.Context, oauthConfig *oauth2.Config, labelChan chan *string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestStateString := ctx.Value(oauthStateStringContextKey).(string)
		responseStateString := r.FormValue("state")
		if responseStateString != requestStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", requestStateString, responseStateString)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		v := url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {code},
			"client_id":    {oauthConfig.ClientID},
			"redirect_uri": {oauthConfig.RedirectURL},
		}
		req, err := http.NewRequest("POST", oauthConfig.Endpoint.TokenURL, strings.NewReader(v.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("oauthoauthConfig.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		//body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("cannot fetch token: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		var tj tokenJSON
		if err = json.Unmarshal(body, &tj); err != nil {
			fmt.Printf("invaid token response: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		token := &oauth2.Token{
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			Expiry:       tj.expiry(),
		}

		//TODO parse token for username
		//store token in native storage
		err = credstore.Set(meroxaLabel, "meroxa-cli", token.AccessToken)
		if err != nil {
			fmt.Printf("Cant store token: %v", err)
		}
		// show success page
		successPage := `
		<div style="height:100px; width:100%!; display:flex; flex-direction: column; justify-content: center; align-items:center; background-color:#2ecc71; color:white; font-size:22"><div>Success!</div></div>
		<p style="margin-top:20px; font-size:18; text-align:center">You are authenticated, you can now return to the program. This will auto-close</p>
		<script>window.onload=function(){setTimeout(this.close, 5000)}</script>
		`
		sm := ""
		fmt.Fprintf(w, successPage)
		// quitSignalChan <- quitSignal
		labelChan <- &sm
	}
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

type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"` // at least PayPal returns string, while most return number
	Expires      expirationTime `json:"expires"`    // broken Facebook spelling of expires_in
}
type expirationTime int32

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	if v := e.Expires; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}
