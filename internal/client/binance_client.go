package client

import "net/http"

type BinanceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewBinanceClient(baseURL string) *BinanceClient {
	return &BinanceClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}
