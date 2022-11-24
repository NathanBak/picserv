package dbpics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// MakeUserAuthenticationURL returns a URL that can be used to login to Dropbox, authorize the app
// (specified by the appKey), and disply a code that can be passed to the GetUserAccessToken()
// function in order to get a refresh token.
func MakeUserAuthenticationURL(appKey string) string {
	return fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&token_access_type=offline&response_type=code", appKey)
}

// GetOfflineAccessToken returns a user offline access token which will have both a short-lived
// access token as well as a refresh token that can be used to obtain fresh access tokens (the
// RefreshOfflineAccessToken() can do that).
func GetOfflineAccessToken(appKey, appSecret, userCode string) (OfflineAccessToken, error) {

	oat := OfflineAccessToken{}

	params := url.Values{}
	params.Add("code", userCode)
	params.Add("grant_type", `authorization_code`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://api.dropbox.com/oauth2/token", body)
	if err != nil {
		return oat, err
	}
	req.SetBasicAuth(appKey, appSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return oat, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return oat, err
	}

	err = json.Unmarshal(buf, &oat)

	return oat, err
}

// OfflineAccessToken contains information that can be used to authenticate API calls to DropBox.
// The AccessToken can be used to make "regular" calls and the RefreshToken can be used to obtain
// new AccessTokens.
type OfflineAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	UID          string `json:"uid"`
	AccountID    string `json:"account_id"`
}

// RefreshOfflineAccessToken uses a refresh token to get a new access token.
func RefreshOfflineAccessToken(appKey, appSecret, userRefreshToken string) (OfflineAccessToken, error) {

	oat := OfflineAccessToken{}

	params := url.Values{}
	params.Add("refresh_token", userRefreshToken)
	params.Add("grant_type", `refresh_token`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://api.dropbox.com/oauth2/token", body)
	if err != nil {
		return oat, err
	}
	req.SetBasicAuth(appKey, appSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return oat, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return oat, err
	}

	err = json.Unmarshal(buf, &oat)

	return oat, err
}
