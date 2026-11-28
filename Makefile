include .env

run:
	@FINNHUB_API_KEY=${FINNHUB_API_KEY} \
	GH_TOKEN=${GH_TOKEN} \
	go run .

test:
	@FINNHUB_API_KEY=${FINNHUB_API_KEY} \
	GH_TOKEN=${GH_TOKEN} \
	go test -v ./

.PHONY: dist
dist:
	@FINNHUB_API_KEY=${FINNHUB_API_KEY} \
	GH_TOKEN=${GH_TOKEN} \
	scripts/dist.sh