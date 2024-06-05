test:
	go test -race -vet=off ./...

run:
	PORT=8080 go run --race cmd/standalone/*.go

run_enricher:
	PORT=8080 go run --race cmd/enricher/*.go

build_standalone:
	go build -o findcastles cmd/standalone/*.go
