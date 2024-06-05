package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/buarki/find-castles/enricher"
	"github.com/buarki/find-castles/htmlfetcher"
	"github.com/buarki/find-castles/httpclient"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT env var")
	}
	httpClient := httpclient.New()
	castlesEnricher := enricher.New(runtime.NumCPU(), httpClient, htmlfetcher.Fetch)

	fs := http.FileServer(http.Dir("./cmd/standalone/public"))
	http.Handle("/", fs)
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		enrichedCastles, enrichmentErrs := castlesEnricher.Enrich(r.Context())

		for {
			select {
			case <-r.Context().Done():
				log.Println("request was canceled")
				return
			case err, ok := <-enrichmentErrs:
				if ok {
					log.Println("received error:", err)
				}
			case castle, ok := <-enrichedCastles:
				if ok {
					cb, err := json.Marshal(castle)
					if err != nil {
						log.Printf("failed to marshal castle [%s]: %v", castle.Name, err)
					}

					if _, err := fmt.Fprintf(w, "data: {\"message\": %s}\n\n", string(cb)); err != nil {
						log.Printf("failed to write to response: %v", err)
					}

					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					} else {
						log.Println("response writer does not support flushing")
					}
				} else {
					fmt.Fprintf(w, "data: {\"finished\":\"finished\"}\n\n")
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					} else {
						log.Println("response writer does not support flushing")
					}
					log.Println("finished processing castles")
					return
				}
			}
		}
	})

	fmt.Println("Server listening on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
