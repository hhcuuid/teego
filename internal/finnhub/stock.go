package finnhub

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"teego.dev/market/internal/github"
)

const (
	DIR_DATA    = github.DATA_DIR
	FILE_STOCKS = DIR_DATA + "/stocks.json"
)

type Stock struct {
	Ticker      string  `json:"ticker"`
	Name        string  `json:"name"`
	Date        string  `json:"date"`
	Industry    string  `json:"industry"`
	MCapBillion float64 `json:"mcapbillion"`
	Beta        float64 `json:"beta"`
	Dividend    float64 `json:"dividend"`
	PeTTM       float64 `json:"pettm"`
	Vol10d90d   float64 `json:"vol10d90d"`
}

func (c *client) ProcessSymbol(ticker, date string) (Stock, error) {
	p2, err := c.Profile2(ticker)
	if err != nil {
		return Stock{}, err
	}
	m, err := c.Metric(ticker)
	if err != nil {
		return Stock{}, err
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s.json", DIR_DATA, ticker), m.Series, 0o644)
	if err != nil {
		return Stock{}, err
	}

	return Stock{
		Ticker:      ticker,
		Name:        p2.Name,
		Date:        date,
		Industry:    p2.Industry,
		MCapBillion: p2.MarketCapM / 1000.0,
		Beta:        m.Metric.Beta,
		Dividend:    m.Metric.Dividend,
		PeTTM:       m.Metric.PeTTM,
		Vol10d90d:   (m.AveTradingVolume10 - m.AveTradingVolume90) / m.AveTradingVolume90 * 100.0,
	}, nil
}

func StoreStocks(stocks map[string]Stock) error {
	threeDaysAgo := time.Now().Add(-time.Hour * 24 * 3).Format("2006-01-02")
	stocksFiltered := make([]Stock, 0, len(stocks))

	for _, v := range stocks {
		if v.Date >= threeDaysAgo {
			stocksFiltered = append(stocksFiltered, v)
		}
	}

	data, err := json.Marshal(stocksFiltered)
	if err != nil {
		return err
	}

	return os.WriteFile(FILE_STOCKS, data, 0o644)
}

func LoadStocks() map[string]Stock {
	profiles := []Stock{}

	file, err := os.Open(FILE_STOCKS)
	if err != nil {
		return map[string]Stock{}
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&profiles); err != nil {
		return map[string]Stock{}
	}

	stocks := make(map[string]Stock, len(profiles))
	for _, s := range profiles {
		stocks[s.Ticker] = s
	}

	return stocks
}
