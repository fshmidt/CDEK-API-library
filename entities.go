package CDEK_API_lib

import "time"

type PriceSending struct {
	TariffCode        int     `json:"tariff_code"`
	TariffName        string  `json:"tariff_name"`
	TariffDescription string  `json:"tariff_description"`
	DeliveryMode      int     `json:"delivery_mode"`
	DeliverySum       float64 `json:"delivery_sum"`
	PeriodMin         int     `json:"period_min"`
	PeriodMax         int     `json:"period_max"`
}

type Client struct {
	Token      string
	TestMode   bool
	APIAddress string
	ExpiresAt  time.Time // время истечения срока действия токена
}

type Location struct {
	Address string `json:"address"`
}

type CalculateRequestBody struct {
	FromLocation Location  `json:"from_location"`
	ToLocation   Location  `json:"to_location"`
	Packages     []Package `json:"packages"`
}

type CreateOrderRequestBody struct {
	Type         int       `json:"type"`
	TariffCode   int       `json:"tariff_code"`
	Value        float64   `json:"value"`
	Threshold    int       `json:"threshold"`
	FromLocation Location  `json:"from_location"`
	ToLocation   Location  `json:"to_location"`
	Recipient    Person    `json:"recipient"`
	Sender       Person    `json:"sender"`
	Packages     []Package `json:"packages"`
}

type Person struct {
	Name    string      `json:"name"`
	Phones  PhoneNumber `json:"phones"`
	Company string      `json:"company"`
}

type PhoneNumber struct {
	Number string `json:"number"`
}

type Package struct {
	Number  string `json:"number"`
	Weight  int    `json:"weight"`
	Length  int    `json:"length"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Comment string `json:"comment"`
}

type OrderCreationResponse struct {
	Entity   Entity         `json:"entity"`
	Requests []OrderRequest `json:"requests"`
}

type Entity struct {
	UUID string `json:"uuid"`
}

type OrderRequest struct {
	RequestUUID string       `json:"request_uuid"`
	Type        string       `json:"type"`
	DateTime    string       `json:"date_time"`
	State       string       `json:"state"`
	Errors      []CdekErrors `json:"errors"`
}

type CdekErrors struct {
	Message string `json:"message"`
}

type CheckResponse struct {
	Entity CheckEntity `json:"entity"`
}

type CheckEntity struct {
	UUID     string   `json:"uuid"`
	Statuses []Status `json:"statuses"`
}

type Status struct {
	Name string `json:"name"`
}
