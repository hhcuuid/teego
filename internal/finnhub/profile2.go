package finnhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FilteredError string

type Profile2 struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	Industry   string  `json:"finnhubIndustry"`
	MarketCapM float64 `json:"marketCapitalization"`
}

func (fe FilteredError) Error() string {
	return string(fe)
}

func (c *client) Profile2(ticker string) (Profile2, error) {
	u := c.JoinPath("profile2")

	q := u.Query()
	q.Set("symbol", ticker)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return Profile2{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return Profile2{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return Profile2{}, fmt.Errorf("status %d", resp.StatusCode)
	}

	p2 := Profile2{}
	if err := json.NewDecoder(resp.Body).Decode(&p2); err != nil {
		return Profile2{}, err
	}

	if p2.Country != "US" {
		return Profile2{}, FilteredError("not a US company")
	}
	if p2.MarketCapM < 500 {
		return Profile2{}, FilteredError("marketCap is too small")
	}
	if p2.Industry == "N/A" {
		return Profile2{}, FilteredError("unknown industry")
	}

	return p2, nil
}
