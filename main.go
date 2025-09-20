package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourname/go-urlfetcher/fetcher"
)

func main() {
	file := flag.String("file", "urls.txt", "file with URLs, one per line")
	workers := flag.Int("workers", 5, "number of concurrent workers")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("open urls file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	jobsChan := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("shutdown requested")
		cancel()
	}()
	results := fetcher.NewWorkerPool(ctx, *workers, jobsChan)

	go func() {
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				break
			case jobsChan <- scanner.Text():
			}
		}
		close(jobsChan)
	}()

	var out []fetcher.Result
	for r := range results {
		log.Printf("got result: %s status=%d len=%d err=%s", r.URL, r.Status, r.Length, r.Err)
		out = append(out, r)
	}
	outF, err := os.Create("results.json")
	if err != nil {
		log.Fatalf("create results.json: %v", err)
	}
	enc := json.NewEncoder(outF)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatalf("encode results: %v", err)
	}
	outF.Close()
	log.Println("done")
}
