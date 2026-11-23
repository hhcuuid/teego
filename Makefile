include .env

run:
	@GH_TOKEN=${GH_TOKEN} \
	go run .

test:
	@GH_TOKEN=${GH_TOKEN} \
	go test ./...