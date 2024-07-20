test:
	go test -race -vet=off ./...

run:
	PORT=8080 go run --race cmd/standalone/*.go

run_enricher:
	PORT=8080 DB_URI="mongodb://localhost:27017/find-castles" ENRICHMENT_TIMEOUT_IN_SECONDS=240 go run --race cmd/enricher/*.go

run_site:
	npm run dev --prefix site

build_standalone:
	go build -o findcastles cmd/standalone/*.go

db_up:
	docker-compose up -d

db_down:
	docker-compose down
