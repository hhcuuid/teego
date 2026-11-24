package main

import (
	"log"
	"os"
	"time"

	"teego.dev/market/internal/finnhub"
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

	fhc, err := finnhub.NewClient(os.Getenv("FINNHUB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = fhc.Symbol()
	if err != nil {
		log.Fatal(err)
	}
}
