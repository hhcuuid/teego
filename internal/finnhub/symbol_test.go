package finnhub

import (
	"os"
	"testing"
)

func Test_finnubClient_Symbol(t *testing.T) {
	fhc, err := NewClient(os.Getenv("FINNHUB_API_KEY"))
	if err != nil {
		t.Fatalf("could not construct client: %v", err)
	}

	t.Run("symbols", func(t *testing.T) {
		got, err := fhc.Symbol()
		if err != nil {
			t.Errorf("Symbol() failed: %v", err)
			return
		}

		if len(got) < 1_000 {
			t.Errorf("Len of Symbol() = %v, want more than 1k", got)
		}
	})
}
