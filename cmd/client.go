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
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/meroxa/meroxa-go"
)

const clientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"
const domain = "auth.meroxa.io"

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

func isDebugEnabled() bool {
	val, ok := os.LookupEnv("MEROXA_DEBUG")
	return flagDebug || (ok && val == "1")
}

func client() (*meroxa.Client, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}

	return meroxa.New(accessToken, VersionString(), isDebugEnabled())
}
