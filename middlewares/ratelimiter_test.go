package middlewares

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TetstRateLimiter(t *testing.T) {
	handler := NewHandlerMock(100 * time.Millisecond)
	middleware := NewRateLimiter(5, handler)

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	concurrentRequests := 6

	var (
		lock              sync.Mutex
		succeededRequests int
		failedRequests    int

		wg sync.WaitGroup
	)

	wg.Add(concurrentRequests)
	for range concurrentRequests {
		go func() {
			defer wg.Done()

			recorder := httptest.NewRecorder()
			middleware.ServeHTTP(recorder, request)

			lock.Lock()
			defer lock.Unlock()

			switch recorder.Code {
			case 200:
				succeededRequests++
			case 429:
				failedRequests++
			}
		}()
	}
	wg.Wait()

	if expectedSucceeds := 5; succeededRequests != expectedSucceeds {
		t.Errorf("expected %d succeeds but got %d", expectedSucceeds, succeededRequests)
	}

	if expectedFailures := 1; failedRequests != expectedFailures {
		t.Errorf("expected %d failures but got %d", expectedFailures, failedRequests)
	}
}

func NewHandlerMock(delay time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)

		w.WriteHeader(http.StatusOK)
	})
}
