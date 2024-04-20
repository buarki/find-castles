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

		defer func() {
			close(collectResults)
		}()

		wg.Add(2)
		go collector.CollectForPotugal(r.Context(), httpClient, &wg, collectResults)
		go collector.CollectForUk(r.Context(), httpClient, &wg, collectResults)

		go func() {
			wg.Wait()
			fmt.Fprintf(w, "data: {\"finished\":\"finished\"}\n\n")
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			} else {
				fmt.Println("response writer does not support flushing")
			}
		}()

		for result := range collectResults {
			if result.Err != nil {
				fmt.Printf("failed to collect info for castle [%s], got %v\n", result.Castle.Name, result.Err)
			} else {
				fmt.Printf("%v\n", result.Castle)
				cb, err := json.Marshal(result.Castle)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if cn, ok := w.(http.CloseNotifier); ok {
					select {
					case <-cn.CloseNotify():
						fmt.Println("Client disconnected. Stopping.")
						return
					default: // Client still connected, continue processing
					}
				}

				if _, err := fmt.Fprintf(w, "data: {\"message\": %s}\n\n", string(cb)); err != nil {
					fmt.Printf("failed to write to response, got %v\n", err)
					continue
				}

				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				} else {
					fmt.Println("response writer does not support flushing")
				}
			}
		}
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
