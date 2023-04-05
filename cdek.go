package CDEK_API_lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func NewClient(username string, password string, testMode bool, apiAddr string) (*Client, error) {

	token, expiresAt, err := GetToken(username, password)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s.json", username)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	client := &Client{
		Token:      token,
		TestMode:   testMode,
		APIAddress: apiAddr,
		ExpiresAt:  expiresAt,
	}

	clientData := struct {
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}{
		Username:  username,
		Password:  password,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err = json.NewEncoder(file).Encode(&clientData)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ClientIsExist(username string, password string, testMode bool, apiAddr string) (*Client, bool) {
	filename := fmt.Sprintf("%s.json", username)

	file, err := os.Open(filename)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	var clientData struct {
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	err = json.NewDecoder(file).Decode(&clientData)
	if err != nil {
		return nil, false
	}

	if clientData.ExpiresAt.Before(time.Now()) || clientData.Username != username || clientData.Password != password {
		if clientData.ExpiresAt.Before(time.Now()) {
			err := os.Remove(filename)
			if err != nil {
				log.Fatal("Can't delete expired client.", err)
			}
		}
		return nil, false
	}
	fmt.Println("Client is already exist. No need to get a new token. It'll be expired in", clientData.ExpiresAt.Sub(time.Now()).Seconds(), "seconds.")
	return &Client{
		Token:      clientData.Token,
		TestMode:   testMode,
		APIAddress: apiAddr,
		ExpiresAt:  clientData.ExpiresAt,
	}, true
}

func (c *Client) ValidateAddress(address string) (res bool, uniformAddr string, err error) {

	previousAPI := c.APIAddress
	c.APIAddress = "https://api.edu.cdek.ru/v2/calculator/tarifflist"

	if _, err := c.Calculate(address, address, Package{"0", 1, 1, 1, 1, ""}); err != nil {
		c.APIAddress = previousAPI
		return false, "", err
	}
	c.APIAddress = previousAPI
	return true, address, nil
}

func (c *Client) GetStatus(orderID string) (status string, err error) {

	query := url.Values{}
	query.Set("authLogin", c.Token)
	url := c.APIAddress + orderID

	fmt.Println("URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CheckResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if len(response.Entity.Statuses) == 0 {
		return "", fmt.Errorf("no statuses found for order with uuid %s", orderID)
	}

	status = response.Entity.Statuses[len(response.Entity.Statuses)-1].Name

	return status, nil
}

func (c *Client) CreateOrder(addrFrom string, addrTo string, size Package, typeSending int) (orderID string, err error) {

	if ok, _, _ := c.ValidateAddress(addrFrom); ok == false {
		return "", errors.New("bad addrFrom: " + addrFrom)
	}
	if ok, _, _ := c.ValidateAddress(addrTo); ok == false {
		return "", errors.New("bad addrTo")
	}

	body := CreateOrderRequestBody{
		Type:       2,
		TariffCode: typeSending,
		Value:      0.0,
		Threshold:  0,
		FromLocation: Location{
			Address: addrFrom,
		},
		ToLocation: Location{
			Address: addrTo,
		},
		Recipient: Person{
			Name: "Вася Пупкин",
			Phones: PhoneNumber{
				"79854595959",
			},
		},
		Sender: Person{
			Name: "Пася Вупкин",
			Phones: PhoneNumber{
				"79854595958",
			},
			Company: "Roga & Kopita",
		},
		Packages: []Package{
			size,
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("json", "true")

	query.Set("authLogin", c.Token)
	url := c.APIAddress

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Добавляем заголовок Authorization с bearer-token авторизации из структуры Client.
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result OrderCreationResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Requests) != 0 && result.Requests[0].State != "ACCEPTED" {

		return "", errors.New(result.Requests[0].Errors[0].Message)
	}
	return result.Entity.UUID, nil

}

func (c *Client) Calculate(addrFrom string, addrTo string, size Package) ([]PriceSending, error) {

	body := CalculateRequestBody{
		FromLocation: Location{
			Address: addrFrom,
		},
		ToLocation: Location{
			Address: addrTo,
		},
		Packages: []Package{
			{
				Number: "0",
				Weight: size.Weight,
				Length: size.Length,
				Width:  size.Width,
				Height: size.Height,
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("json", "true")

	if c.TestMode {
		query.Set("test", "1")
	}
	query.Set("authLogin", c.Token)
	url := c.APIAddress

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Добавляем заголовок Authorization с bearer-token авторизации из структуры Client.
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		TariffCodes []PriceSending `json:"tariff_codes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.TariffCodes) == 0 {
		return nil, errors.New("Bad address")
	}
	return result.TariffCodes, nil
}
