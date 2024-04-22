package handlers

import (
	"net/http"
	"net/http/httptest"
	"ringbuffer/counter"
	"testing"
	"time"
)

func TestCounterHandler(t *testing.T) {
	counter := NewCounterMock()
	handler := NewCounterHandler(counter)

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("wanted %v status code but got %v", http.StatusOK, status)
	}

	expected := `requests received during the previous 60 seconds: 1`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}

	if ttlCalledCount := counter.TTLCalledCount; ttlCalledCount != 1 {
		t.Errorf("expected counter TTL func to be called only once, but called %v", ttlCalledCount)
	}

	if incCalledCount := counter.IncCalledCount; incCalledCount != 1 {
		t.Errorf("expected counter Inc func to be called only once, but called %v", incCalledCount)
	}

	if countCalledCount := counter.CountCalledCount; countCalledCount != 1 {
		t.Errorf("expected counter Count func to be called only once, but called %v", countCalledCount)
	}
}

type counterMock struct {
	counter.Counter

	TTLCalledCount   int
	IncCalledCount   int
	CountCalledCount int
}

func NewCounterMock() *counterMock {
	return &counterMock{}
}

func (c *counterMock) TTL() int {
	c.TTLCalledCount++

	return 60
}

func (c *counterMock) Inc(at time.Time) {
	c.IncCalledCount++
}

func (c *counterMock) Count() int {
	c.CountCalledCount++

	return c.CountCalledCount
}
