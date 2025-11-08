package client

import (
	"context"
	"encoding/json"
	"fmt"
	"go-finance/internal/models"
	"net/http"
	"net/url"
	"time"
)

type FinmindClient struct {
	HTTPClient *http.Client
	// apiKey     string
	BaseURL string
}

func NewFinmindClient(baseURL string) *FinmindClient {
	return &FinmindClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *FinmindClient) GetTaiwanStockPrice(ctx context.Context, dataset, dataID, startDate string) (*models.FinmindResponse, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	// u.Host would look like "api.finmindtrade.com"
	u.Path = "/api/v4/data"
	q := u.Query()
	q.Set("dataset", dataset)
	q.Set("data_id", dataID)
	q.Set("start_date", startDate)
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("finmind: status %d", res.StatusCode)
	}

	var r models.FinmindResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
