package handlers

import (
	"fmt"
	"net/http"
	"ringbuffer/counter"
	"time"
)

type counterHandler struct {
	counter     counter.Counter
	timeToSleep time.Duration
}

var _ http.Handler = &counterHandler{}

func NewCounterHandler(counter counter.Counter, timeToSleep time.Duration) *counterHandler {
	return &counterHandler{
		counter:     counter,
		timeToSleep: timeToSleep,
	}
}

func (c *counterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(c.timeToSleep)

	c.counter.Inc(time.Now())

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "requests received during the previous %d seconds: %d", c.counter.TTL(), c.counter.Count())
}
