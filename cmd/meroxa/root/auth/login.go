/*
Copyright © 2022 Meroxa Inc

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

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"

	"github.com/fatih/color"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
)

var (
	_ builder.CommandWithDocs    = (*Login)(nil)
	_ builder.CommandWithLogger  = (*Login)(nil)
	_ builder.CommandWithExecute = (*Login)(nil)
	_ builder.CommandWithConfig  = (*Login)(nil)
)

type Login struct {
	logger log.Logger
	config config.Config
}

func (l *Login) Usage() string {
	return "login"
}

func (l *Login) Docs() builder.Docs {
	return builder.Docs{
		Short: "Login or Sign up to the Meroxa Platform",
	}
}

func (l *Login) Execute(ctx context.Context) error {
	// initialize the code verifier
	err := l.login(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (l *Login) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *Login) Config(cfg config.Config) {
	l.config = cfg
}

// AuthorizeUser implements the PKCE OAuth2 flow.
func (l *Login) authorizeUser(ctx context.Context, clientID, authDomain, audience, redirectURL string) {
	// initialize the code verifier
	var CodeVerifier, _ = cv.CreateCodeVerifier()

	// Create code_challenge with S256 method
	codeChallenge := CodeVerifier.CodeChallengeS256()

	// construct the authorization URL (with Auth0 as the authorization provider)
	authorizationURL := fmt.Sprintf(
		"https://%s/authorize?audience=%s"+
			`&scope=openid%%20email%%20offline_access%%20user`+
			"&response_type=code&client_id=%s"+
			"&code_challenge=%s"+
			"&code_challenge_method=S256&redirect_uri=%s",
		authDomain, audience, clientID, codeChallenge, redirectURL)
	l.logger.Infof(ctx, color.CyanString(authorizationURL))

	// start a web server to listen on a callback URL
	server := &http.Server{Addr: redirectURL, ReadHeaderTimeout: time.Minute}

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		// get the authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			l.logger.Errorf(ctx, "meroxa: Url Param 'code' is missing")
			_, _ = io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			l.cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		codeVerifier := CodeVerifier.String()
		accessToken, refreshToken, err := l.getAccessTokenAuth(r.Context(), clientID, codeVerifier, code, redirectURL)
		if err != nil {
			l.logger.Errorf(ctx, "meroxa: could not get access token")
			_, _ = io.WriteString(w, "Error: could not retrieve tokens\n")

			// close the HTTP server and return
			l.cleanup(server)
			return
		}

		l.config.Set(global.AccessTokenEnv, accessToken)
		l.config.Set(global.RefreshTokenEnv, refreshToken)

		// return an indication of success to the caller
		_, _ = io.WriteString(w, `
			<html>
				<div style="width:100%!; color:#282D39; display:flex; flex-direction: column; justify-content:center; 
						align-items:center; margin-top:40px;">
					<img src="https://meroxa-public-assets.s3.us-east-2.amazonaws.com/MeroxaTransparent%402x.png" alt="Meroxa logo"
						 width="150" padding="2000px">
					<h1 style="margin-top:40px; font-size:43px; text-align:left; color: #282D39; font-family: Arial; font-weight: bold;">
						Successfully logged in
					</h1>
					<p style="margin-top:17px; font-size:18px; text-align:left; color: #282D39; font-family: Arial;">
						You can close this window now and return to the CLI.
					</p>
				</div>
			</html>`)

		l.logger.Infof(ctx, "Successfully logged in.")

		// close the HTTP server
		l.cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		l.logger.Errorf(ctx, "meroxa: bad redirect URL: %s\n", err)
		os.Exit(1)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	listener, err := net.Listen("tcp", port)
	if err != nil {
		l.logger.Errorf(ctx, "meroxa: can't listen to port %s: %s\n", port, err)
		os.Exit(1)
	}

	// open a browser window to the authorizationURL
	err = open.Start(authorizationURL)
	if err != nil {
		l.logger.Errorf(ctx, "meroxa: can't open browser to URL %s: %s\n", authorizationURL, err)
		os.Exit(1)
	}

	// start the blocking web server loop
	// this will exit when the handler gets fired and calls server.Close()
	_ = server.Serve(listener)
}

func (l *Login) login(ctx context.Context) error {
	l.logger.Infof(ctx, color.CyanString("You will now be taken to your browser for authentication or open the url below in a browser."))
	l.logger.Infof(ctx, global.GetMeroxaAuthCallbackURL())
	l.authorizeUser(
		ctx,
		global.GetMeroxaAuthClientID(),
		global.GetMeroxaAuthDomain(),
		global.GetMeroxaAuthAudience(),
		global.GetMeroxaAuthCallbackURL(),
	)
	return nil
}

// cleanup closes the HTTP server.
func (l *Login) cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	l.config.Set(global.UserAccountUUID, "")
	go server.Close()
}

// getAccessTokenAuth trades the authorization code retrieved from the first OAuth2 leg for an access token.
func (l *Login) getAccessTokenAuth(
	ctx context.Context,
	clientID, codeVerifier, authorizationCode, callbackURL string,
) (accessToken, refreshToken string, err error) {
	// set the url and form-encoded data for the POST to the access token endpoint
	// this URL should actually be taken from meroxa.OAuth2Endpoint
	tokenURL := fmt.Sprintf("https://%s/oauth/token", global.GetMeroxaAuthDomain())
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, payload)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		l.logger.Errorf(ctx, "meroxa: HTTP error: %s", err)
		return "", "", err
	}

	// process the response
	defer res.Body.Close()
	var responseData map[string]interface{}
	body, _ := io.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		l.logger.Errorf(ctx, "meroxa: JSON error: %s", err)
		return "", "", err
	}

	// retrieve the access token out of the map, and return to caller
	accessToken = responseData["access_token"].(string)
	refreshToken = responseData["refresh_token"].(string)
	return accessToken, refreshToken, nil
}
