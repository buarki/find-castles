test:
	go test -race -vet=off ./...

run:
	PORT=8080 go run --race cmd/*.go

build:
	go build -o findcastles cmd/*.go
