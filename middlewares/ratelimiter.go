package middlewares

import (
	"fmt"
	"net/http"
	"sync"
)

type rateLimiter struct {
	next      http.Handler
	threshold int

	indicator int
	lock      sync.RWMutex
}

var _ http.Handler = &rateLimiter{}

func NewRateLimiter(threshold int, next http.Handler) *rateLimiter {
	return &rateLimiter{
		threshold: threshold,
		next:      next,
	}
}

func (m *rateLimiter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if m.thresholdExceeded() {
		rw.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(rw, "you have exceeded %d parallel requests threshold", m.threshold)

		return
	}

	m.inc()
	m.next.ServeHTTP(rw, r)
	m.dec()
}

func (m *rateLimiter) inc() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.indicator++
}

func (m *rateLimiter) dec() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.indicator--
}

func (m *rateLimiter) thresholdExceeded() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.indicator >= m.threshold
}
