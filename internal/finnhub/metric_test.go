package finnhub

import (
	"fmt"
	"os"
	"testing"
)

func Test_finnubClient_Metric(t *testing.T) {

	fhClient, err := NewClient(os.Getenv("FINNHUB_API_KEY"))
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

			got, gotErr := fhClient.Metric(tt.symbol)
			if gotErr != nil {
				t.Errorf("Metric() failed: %v", gotErr)
				return
			}

			if got.PeTTM == 0.0 {
				t.Errorf("ForwardPE = %v, want != 0.0", got.PeTTM)
			}
			if len(got.Series) < 1000 {
				t.Errorf("Len(Series) = %v, want > 1000", len(got.Series))
			}
			fmt.Println(got.Series.String())
		})
	}
}
