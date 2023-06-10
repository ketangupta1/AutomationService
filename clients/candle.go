package clients

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/TredingInGo/AutomationService/utils"
)

const (
	name           string        = "smartapi-go"
	requestTimeout time.Duration = 7000 * time.Millisecond
	baseURI        string        = "https://apiconnect.angelbroking.com/"
	historyApiKey  string        = "MN9K2rhC"
)

type HistoricClient struct {
	clientCode  string
	password    string
	accessToken string
	debug       bool
	baseURI     string
	apiKey      string
	httpClient  HTTPClient
}

func NewHistoricClient(clientCode string, password string) HistoricClient {
	client := HistoricClient{
		clientCode:  clientCode,
		password:    password,
		accessToken: "",
		debug:       false,
		baseURI:     baseURI,
		apiKey:      historyApiKey,
		httpClient:  nil,
	}

	// Create a default http handler with default timeout.
	client.httpClient = NewHTTPClient(&http.Client{
		Timeout:   requestTimeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}, nil, client.debug)

	return client
}

func (hc *HistoricClient) doEnvelope(method, uri string, params map[string]interface{}, headers http.Header, authorization ...bool) ([]byte, error) {
	// Send custom headers set
	if headers == nil {
		headers = map[string][]string{}
	}

	headers.Add("X-PrivateKey", hc.apiKey)
	if authorization != nil && authorization[0] {
		headers.Add("Authorization", "Bearer "+hc.accessToken)
	}

	return hc.httpClient.DoEnvelope(method, hc.baseURI+uri, params, headers)
}

// GenerateSession totp used is required for 2 factor authentication
func (hc *HistoricClient) GenerateSession(totp string) (UserSession, error) {

	// construct url values
	params := make(map[string]interface{})
	params["clientcode"] = hc.clientCode
	params["password"] = hc.password
	params["totp"] = totp

	var session UserSession
	resp, err := hc.doEnvelope(http.MethodPost, utils.URILogin, params, nil)
	// Set accessToken on successful session retrieve
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(resp, &session)
	if err != nil {
		return session, err
	}

	hc.SetAccessToken(session.AccessToken)

	return session, err
}

// SetAccessToken sets the access token to the Kite Connect instance.
func (hc *HistoricClient) SetAccessToken(accessToken string) {
	hc.accessToken = accessToken
}

// CandleParams represents parameters for getting CandleData.
type CandleParams struct {
	Exchange    string `json:"exchange"`
	SymbolToken string `json:"symboltoken"`
	Interval    string `json:"interval"`
	FromDate    string `json:"fromdate"`
	ToDate      string `json:"todate"`
}

func (c CandleParams) getParams() map[string]interface{} {
	params := make(map[string]interface{})

	if c.Exchange != "" {
		params["exchange"] = c.Exchange
	}

	if c.SymbolToken != "" {
		params["symboltoken"] = c.Exchange
	}

	if c.Interval != "" {
		params["interval"] = c.Interval
	}

	if c.FromDate != "" {
		params["fromdate"] = c.FromDate
	}

	if c.ToDate != "" {
		params["todate"] = c.ToDate
	}

	return params
}

type CandleResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int       `json:"volume"`
}

func (hc *HistoricClient) GetCandleData(candleParam CandleParams) ([]CandleResponse, error) {
	var candleData []CandleResponse
	params := candleParam.getParams()
	resp, err := hc.doEnvelope(http.MethodPost, utils.URICandleData, params, nil, true)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &candleData)
	if err != nil {
		return nil, err
	}

	return candleData, err
}
