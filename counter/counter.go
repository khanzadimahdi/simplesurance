package counter

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

// Counter interface
type Counter interface {
	TTL() int
	Inc(at time.Time)
	Count() int
	Store(writer io.Writer) error
	Load(reader io.Reader) error
}

// counter implements the sliding-window pattern to count number of Incs in within given TTL
type counter struct {
	ttl int

	// ring buffer
	lock      sync.Mutex
	Buf       []int       `json:"buf"`
	OccuredAt []time.Time `json:"occured_at"`
}

var _ Counter = &counter{}

// NewCounter returns a new instance of counter
func NewCounter(ttl int) *counter {
	return &counter{
		ttl:       ttl,
		Buf:       make([]int, ttl),
		OccuredAt: make([]time.Time, ttl),
	}
}

// TTL returns Time to Live of lower bound of the sliding window
func (c *counter) TTL() int {
	return c.ttl
}

// Inc increases the counter
func (c *counter) Inc(at time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	head := at.Unix() % int64(c.ttl)
	if at.Sub(c.OccuredAt[head]).Seconds() >= 1 {
		c.Buf[head] = 0
	}

	c.Buf[head]++
	c.OccuredAt[head] = at
}

// Count returns count of Inc calls within given defined TTL
func (c *counter) Count() int {
	var total int
	now := time.Now()

	c.lock.Lock()
	defer c.lock.Unlock()

	for i := range c.Buf {
		if now.Sub(c.OccuredAt[i]).Seconds() <= float64(c.ttl) {
			total += c.Buf[i]
		}
	}

	return total
}

// Store writes the inner state of counter to given io.writer
func (c *counter) Store(writer io.Writer) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return json.NewEncoder(writer).Encode(c)
}

// Load loads the inner state of counter from given io.reader.
func (c *counter) Load(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(c)
}
