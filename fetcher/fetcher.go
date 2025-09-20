package fetcher

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
	Length int    `json:"length"`
	Err    string `json:"error,omitempty"`
}

func NewWorkerPool(ctx context.Context, workers int, jobs <-chan string) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 10 * time.Second}

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case u, ok := <-jobs:
				if !ok {
					return
				}
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
				if err != nil {
					results <- Result{URL: u, Err: err.Error()}
					continue
				}
				resp, err := client.Do(req)
				if err != nil {
					results <- Result{URL: u, Err: err.Error()}
					continue
				}
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					results <- Result{URL: u, Status: resp.StatusCode, Err: err.Error()}
					continue
				}
				results <- Result{URL: u, Status: resp.StatusCode, Length: len(body)}
			}
		}
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
