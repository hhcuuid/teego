package main

import (
	"errors"
	"log"
	"os"
	"time"

	"teego.dev/market/internal/finnhub"
	"teego.dev/market/internal/github"
)

func main() {
	log.Default().SetFlags(
		log.Ltime |
			log.LUTC |
			log.Lshortfile,
	)
	log.Print("@start")
	start := time.Now()
	defer func() {
		log.Printf("@stop [%v]", time.Since(start))
	}()

	if err := github.CleanFiles(); err != nil {
		log.Fatal(err)
	}
	if err := github.CloneMarketRepo(); err != nil {
		log.Fatal(err)
	}
	defer github.CleanFiles()

	fhc, err := finnhub.NewClient(os.Getenv("FINNHUB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	symbols, err := fhc.Symbol()
	if err != nil {
		log.Fatal(err)
	}

	stocks := finnhub.LoadStocks()
	today := time.Now().Format("2006-01-02")
	filteredError := new(finnhub.FilteredError)

	for i, sym := range symbols {
		if i%128 == 0 {
			log.Printf("progress [ %04d / %d ]", i, len(symbols))
		}
		stock, err := fhc.ProcessSymbol(sym, today)
		if err != nil {
			if !errors.As(err, filteredError) {
				log.Println(err)
			}
			continue
		}
		stocks[sym] = stock
	}

	if err := finnhub.StoreStocks(stocks); err != nil {
		log.Fatal(err)
	}
	if err := github.AddCommitPush(); err != nil {
		log.Fatal(err)
	}
}
