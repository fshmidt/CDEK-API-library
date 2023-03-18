package CDEK_API_lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func GetToken(account string, password string) (string, error) {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", account)
	values.Set("client_secret", password)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.edu.cdek.ru/v2/oauth/token", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to get token. status code: %d", resp.StatusCode)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
