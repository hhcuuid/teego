package finnhub

import (
	"os"
	"testing"
)

func Test_finnubClient_Profile2(t *testing.T) {

	fhc, err := NewClient(os.Getenv("FINNHUB_API_KEY"))
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	tests := []struct {
		name   string
		symbol string
	}{
		{
			name:   "Apple Inc",
			symbol: "AAPL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p2, err := fhc.Profile2(tt.symbol)
			if err != nil {
				t.Errorf("Profile2() failed: %v", err)
				return
			}

			if p2.MarketCapM < 1000 {
				t.Errorf("Profile2() = %v, want > 1", p2.MarketCapM)
			}
		})
	}
}
