package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWorkerPoolSimple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	jobs := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := NewWorkerPool(ctx, 2, jobs)

	go func() {
		jobs <- ts.URL
		close(jobs)
	}()

	got := <-results
	if got.Status != http.StatusOK {
		t.Fatalf("expected 200 got %d", got.Status)
	}
}
