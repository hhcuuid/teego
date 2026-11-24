package main

import (
	"errors"
	"os"
	"testing"
	"time"

	"teego.dev/market/internal/finnhub"
	"teego.dev/market/internal/github"
)

func TestFullWorkflow(t *testing.T) {
	if err := github.CleanFiles(); err != nil {
		t.Fatalf("CleanFiles() failed: %v", err)
	}
	if err := github.CloneMarketRepo(); err != nil {
		t.Fatalf("CloneMarketRepo() failed: %v", err)
	}
	defer github.CleanFiles()

	fhc, err := finnhub.NewClient(os.Getenv("FINNHUB_API_KEY"))
	if err != nil {
		t.Fatalf("could not construct client: %v", err)
	}

	symbols, err := fhc.Symbol()
	if err != nil {
		t.Fatalf("could not retrieve symbols: %v", err)
	}

	stocks := finnhub.LoadStocks()
	today := time.Now().Format("2006-01-02")
	filteredError := new(finnhub.FilteredError)

	for i, sym := range symbols {
		if i == 10 {
			break
		}
		stock, err := fhc.ProcessSymbol(sym, today)
		if err != nil {
			if !errors.As(err, filteredError) {
				t.Fatalf("could not process symbol %s: %v", sym, err)
			}
			continue
		}
		stocks[sym] = stock
	}

	if err := finnhub.StoreStocks(stocks); err != nil {
		t.Fatalf("could not store stocks: %v", err)
	}

	if err := github.AddCommitPush(); err != nil {
		t.Fatalf("AddCommitPush() failed: %v", err)
	}
}
