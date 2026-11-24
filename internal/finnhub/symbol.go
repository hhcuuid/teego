package finnhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Symbol struct {
	Symbol string `json:"symbol"`
	Mic    string `json:"mic"`
}

/*
get all symbols with
exchange=US,
currency=USD,
securityType="Common Stock"
https://finnhub.io/docs/api#stock-symbols
*/
func (c *client) Symbol() ([]string, error) {
	u := c.JoinPath("symbol")

	q := u.Query()
	q.Set("exchange", "US")
	q.Set("currency", "USD")
	q.Set("securityType", "Common Stock") // ETP, REIT
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var symbols = make([]string, 0, 5120)
	dec := json.NewDecoder(resp.Body)

	// 1. Read the opening bracket of the JSON array '['
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if _, ok := t.(json.Delim); !ok {
		return nil, fmt.Errorf("expected JSON array start")
	}

	var s Symbol
	for dec.More() {
		if err := dec.Decode(&s); err != nil {
			return nil, err
		}

		if s.Mic != "OOTC" && !strings.Contains(s.Symbol, ".") {
			symbols = append(symbols, s.Symbol)
		}
	}

	t, err = dec.Token()
	if err != nil {
		return nil, err
	}
	if _, ok := t.(json.Delim); !ok {
		return nil, fmt.Errorf("expected JSON array end")
	}

	return symbols, nil
}
