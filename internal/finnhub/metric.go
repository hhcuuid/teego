package finnhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Metric struct {
	Beta               float64 `json:"beta"`
	Dividend           float64 `json:"currentDividendYieldTTM"`
	PeTTM              float64 `json:"peTTM"`
	AveTradingVolume10 float64 `json:"10DayAverageTradingVolume"`
	AveTradingVolume90 float64 `json:"3MonthAverageTradingVolume"`
}

type MetricResponse struct {
	Metric `json:"metric"`
	Series json.RawMessage `json:"series"`
}

func (c *client) Metric(ticker string) (MetricResponse, error) {
	u := c.JoinPath("metric")

	q := u.Query()
	q.Set("symbol", ticker)
	q.Set("metric", "all")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return MetricResponse{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return MetricResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return MetricResponse{}, fmt.Errorf("status %d", resp.StatusCode)
	}

	mr := MetricResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return MetricResponse{}, err
	}

	return mr, nil
}
