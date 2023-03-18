package CDEK_API_lib

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type PriceSending struct {
	TariffCode        int     `json:"tariff_code"`
	TariffName        string  `json:"tariff_name"`
	TariffDescription string  `json:"tariff_description"`
	DeliveryMode      int     `json:"delivery_mode"`
	DeliverySum       float64 `json:"delivery_sum"`
	PeriodMin         int     `json:"period_min"`
	PeriodMax         int     `json:"period_max"`
}

type Size struct {
	Width  int `json:"width"`
	Length int `json:"length"`
	Height int `json:"height"`
	Weight int `json:"weight"`
}

type Client struct {
	Token      string
	TestMode   bool
	APIAddress string
}

func NewClient(username string, password string, testMode bool, apiAddr string) (*Client, error) {
	token, err := GetToken(username, password)
	if err != nil {
		return nil, err
	}

	return &Client{
		Token:      token,
		TestMode:   testMode,
		APIAddress: apiAddr,
	}, nil
}

func (c *Client) Calculate(addrFrom string, addrTo string, size Size) ([]PriceSending, error) {
	body := map[string]interface{}{
		"from_location": map[string]interface{}{
			"address": addrFrom,
		},
		"to_location": map[string]interface{}{
			"address": addrTo,
		},
		"packages": []map[string]interface{}{
			{
				"weight": size.Weight,
				"length": size.Length,
				"width":  size.Width,
				"height": size.Height,
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

	return result.TariffCodes, nil
}
