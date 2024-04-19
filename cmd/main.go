package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/buarki/find-castles/collector"
	"github.com/buarki/find-castles/htttpclient"
)

func main() {
	httpClient := htttpclient.New()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		var wg sync.WaitGroup
		collectResults := make(chan collector.CollectResult)

		wg.Add(1)
		go collector.CollectForPotugal(r.Context(), httpClient, &wg, collectResults)

		wg.Add(1)
		go func() {
			defer wg.Done()

			fmt.Println("waiting...")
			for result := range collectResults {
				if result.Err != nil {
					fmt.Printf("failed to collect info for castle [%s], got %v\n", result.Castle.Name, result.Err)
				} else {
					fmt.Printf("%v\n", result.Castle)
					cb, err := json.Marshal(result.Castle)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Fprintf(w, "data: {\"message\": %s}\n\n", string(cb))
					w.(http.Flusher).Flush()
				}
			}
			close(collectResults)
			fmt.Println("closed...")
		}()
		wg.Wait()
		fmt.Println(">>>> FINISHED <<<")
		fmt.Fprintf(w, "data: {\"finished\":\"finished\"}\n\n")
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
