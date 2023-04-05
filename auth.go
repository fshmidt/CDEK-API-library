package CDEK_API_lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetToken(account string, password string) (string, time.Time, error) {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", account)
	values.Set("client_secret", password)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.edu.cdek.ru/v2/oauth/token", strings.NewReader(values.Encode()))
	if err != nil {
		return "", time.Time{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("unable to get token. status code: %d", resp.StatusCode)
	}

	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(response.ExpiresIn))

	return response.AccessToken, expiresAt, nil
}
