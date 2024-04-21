test:
	go test -race -vet=off ./...

run:
	go run --race cmd/*.go

build:
	go build cmd/
