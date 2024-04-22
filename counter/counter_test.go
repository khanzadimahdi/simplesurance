package counter

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestCounter_IncAndCount(t *testing.T) {
	testCases := []struct {
		TTL      int
		Duration int
		Attempts int
		Count    int
	}{
		{TTL: 10, Duration: 30, Attempts: 10, Count: 100},
		{TTL: 60, Duration: 75, Attempts: 10, Count: 600},
		{TTL: 75, Duration: 75, Attempts: 10, Count: 750},
		{TTL: 200, Duration: 75, Attempts: 10, Count: 750},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			t.Parallel()

			gap := 1 // 1 second gap to be able to do the count before TTL
			counter := NewCounter(testCase.TTL + gap)

			now := time.Now().Add(-time.Duration(testCase.Duration) * time.Second)
			for i := 0; i < testCase.Duration; i++ {
				at := now.Add(time.Duration(i) * time.Second)
				for j := 0; j < testCase.Attempts; j++ {
					counter.Inc(at)
				}
			}

			if got := counter.Count(); testCase.Count != got {
				t.Errorf("expected count to be %d but got %d", testCase.Count, got)
			}
		})
	}
}

func TestCounter_StoreAndLoad(t *testing.T) {
	ttl := 15
	counter := NewCounter(ttl)

	items := []struct {
		Attempts int
		At       time.Time
	}{
		{Attempts: 5, At: time.Now().Add(-6 * time.Second)},
		{Attempts: 4, At: time.Now().Add(-2 * time.Second)},
		{Attempts: 3, At: time.Now()},
	}

	for _, item := range items {
		for j := 0; j < item.Attempts; j++ {
			counter.Inc(item.At)
		}
	}

	var buf bytes.Buffer

	if err := counter.Store(&buf); err != nil {
		t.Fatal("unexpected error", err)
	}

	if err := counter.Load(&buf); err != nil {
		t.Fatal("unexpected error", err)
	}

	want := 12
	got := counter.Count()
	if counter.Count() != 12 {
		t.Fatalf("expected count to be %d but got %d", want, got)
	}
}

func BenchmarkTTL60(b *testing.B) {
	counter := NewCounter(60)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		at := time.Now().Add(time.Duration(n) * time.Second)
		b.StartTimer()

		counter.Inc(at)
		counter.Count()
	}
}

func BenchmarkIncTTL100(b *testing.B) {
	counter := NewCounter(100)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		at := time.Now().Add(time.Duration(n) * time.Second)
		b.StartTimer()

		counter.Inc(at)
		counter.Count()
	}
}
