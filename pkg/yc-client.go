package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"yc-backend/utils"

	"github.com/samber/lo"
)

type Channel struct {
	ID                      string    `json:"id"`
	Max                     float64   `json:"max"`
	Currency                string    `json:"currency"`
	CountryCurrency         string    `json:"countryCurrency"`
	Status                  string    `json:"status"`
	FeeLocal                int       `json:"feeLocal"`
	CreatedAt               time.Time `json:"createdAt"`
	VendorID                string    `json:"vendorId"`
	Country                 string    `json:"country"`
	FeeUSD                  float64   `json:"feeUSD"`
	Min                     float64   `json:"min"`
	ChannelType             string    `json:"channelType"`
	RampType                string    `json:"rampType"`
	UpdatedAt               time.Time `json:"updatedAt"`
	APIStatus               string    `json:"apiStatus"`
	SettlementType          string    `json:"settlementType"`
	EstimatedSettlementTime int       `json:"estimatedSettlementTime"`
	Balancer                struct{}  `json:"balancer"`
}

type Network struct {
	Code                     string    `json:"code"`
	UpdatedAt                time.Time `json:"updatedAt"`
	Status                   string    `json:"status"`
	ChannelIds               []string  `json:"channelIds"`
	AccountNumberType        string    `json:"accountNumberType"`
	CreatedAt                time.Time `json:"createdAt"`
	ID                       string    `json:"id"`
	Country                  string    `json:"country"`
	Name                     string    `json:"name"`
	CountryAccountNumberType string    `json:"countryAccountNumberType"`
}

type AccountDetail struct {
	Available    float64 `json:"available"`
	Currency     string  `json:"currency"`
	CurrencyType string  `json:"currencyType"`
}

type Rate struct {
	Buy       float64   `json:"buy"`
	Locale    string    `json:"locale"`
	RateID    string    `json:"rateId"`
	Code      string    `json:"code"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ChannelResponse struct {
	Channels []Channel `json:"channels"`
}

type AccountDetailResponse struct {
	AccountDetail []AccountDetail `json:"accounts"`
}

type NetworkResponse struct {
	Networks []Network `json:"networks"`
}

type RateResponse struct {
	Rates []Rate `json:"rates"`
}

// YellowClient struct
type YellowClient struct {
	client                     *http.Client
	baseUrl, apiKey, apiSecret string
}

// NewYellowClient constructor
func NewYellowClient(baseUrl, apiKey, apiSecret string) *YellowClient {
	return &YellowClient{
		client:    &http.Client{},
		baseUrl:   baseUrl,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// httpAuth method to generate authorization headers
func (yc *YellowClient) httpAuth(path, method string, body map[string]interface{}) (map[string]string, error) {
	yc.client.Timeout = time.Second * 10
	yc.client.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
	}
	yc.client.Jar = nil
	yc.client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }

	date := time.Now().UTC().Format(time.RFC3339)
	h := hmac.New(sha256.New, []byte(yc.apiSecret))
	h.Write([]byte(date))
	h.Write([]byte(path))
	h.Write([]byte(method))

	if body != nil && lo.Contains([]string{http.MethodPost, http.MethodPut}, method) {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyHmac := sha256.Sum256(bodyBytes)
		bodyB64 := base64.StdEncoding.EncodeToString(bodyHmac[:])
		h.Write([]byte(bodyB64))
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return map[string]string{
		"X-YC-Timestamp": date,
		"Authorization":  fmt.Sprintf("YcHmacV1 %s:%s", yc.apiKey, signature),
		"Content-Type":   "application/json",
	}, nil
}

// MakeRequest method to make an authorized request
func (yc *YellowClient) MakeRequest(method string, path string, body map[string]interface{}) (*http.Response, error) {
	headers, err := yc.httpAuth(path, method, body)
	if err != nil {
		return nil, err
	}

	log.Printf("header = %v", headers)

	url := yc.baseUrl + path

	var bodyBytes []byte
	if body != nil && lo.Contains([]string{http.MethodPost, http.MethodPut}, method) {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	http.DefaultClient = yc.client
	req, err := http.NewRequest(method, url, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	utils.LoopOverMap(headers, func(k, v string) { req.Header.Set(k, v) })

	resp, err := yc.client.Do(req)
	if err != nil {
		return nil, err
	}

	if !lo.Contains([]int{http.StatusCreated, http.StatusOK, http.StatusNoContent}, resp.StatusCode) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		return resp, fmt.Errorf("failed to submit payment = %v , code = %v", string(body), resp.Status)
	}

	return resp, nil
}

func (yc *YellowClient) GetYellowCardChannels() ([]Channel, error) {
	var channelResponse ChannelResponse
	resp, err := yc.MakeRequest(http.MethodGet, "/business/channels", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &channelResponse)

	if err != nil {
		return nil, err
	}
	return channelResponse.Channels, nil
}

func (yc *YellowClient) GetYellowCardNetworks() ([]Network, error) {
	var networkResponse NetworkResponse
	resp, err := yc.MakeRequest(http.MethodGet, "/business/networks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &networkResponse)

	if err != nil {
		return nil, err
	}

	return networkResponse.Networks, nil
}

func (yc *YellowClient) GetYellowCardRates() ([]Rate, error) {
	var rateResponse RateResponse
	resp, err := yc.MakeRequest(http.MethodGet, "/business/rates", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &rateResponse)

	if err != nil {
		return nil, err
	}
	return rateResponse.Rates, nil
}
