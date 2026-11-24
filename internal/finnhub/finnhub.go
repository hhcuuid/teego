package finnhub

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	hc *http.Client
	*rate.Limiter
	*url.URL
	apiKey string
}

func NewClient(apiKey string) (*client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("empty apiKey")
	}

	base, err := url.Parse("https://finnhub.io/api/v1/stock")
	if err != nil {
		return nil, err
	}

	return &client{
		&http.Client{
			Timeout: time.Second * 10,
		},
		rate.NewLimiter(rate.Every(time.Minute/58), 1),
		base,
		apiKey,
	}, nil
}

func (fc *client) Do(req *http.Request) (*http.Response, error) {
	if err := fc.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("rate limiter cancelled: %w", err)
	}
	req.Header.Set("X-Finnhub-Token", fc.apiKey)
	req.Header.Set("Accept", "application/json")

	return fc.hc.Do(req)
}
