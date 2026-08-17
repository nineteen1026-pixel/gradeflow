// Package tally keeps running point totals for graded submissions.
//
// A Counter is shared by every grading worker in the service, so all of its
// operations are safe to call from multiple goroutines.
package tally

import (
	"sort"
	"sync"
)

// entry is the mutable cell stored for a single key. Keeping the value behind
// a pointer lets the map keep a stable address for a key across its lifetime.
type entry struct {
	n int
}

// Counter accumulates integer point totals keyed by assignment or student id.
type Counter struct {
	mu sync.Mutex
	m  map[string]*entry
}

// New returns an empty Counter ready for use.
func New() *Counter {
	return &Counter{m: make(map[string]*entry)}
}

// Bump adds delta to the total stored under key and returns the new total.
// Keys are created on first use.
func (c *Counter) Bump(key string, delta int) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.m[key]
	if !ok {
		e = &entry{}
		c.m[key] = e
	}

	e.n += delta
	return e.n
}

// Value reports the total currently stored under key, or zero when the key has
// never been bumped.
func (c *Counter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.m[key]
	if !ok {
		return 0
	}
	return e.n
}

// Total sums the totals of every key held by the counter.
func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	sum := 0
	for _, e := range c.m {
		sum += e.n
	}
	return sum
}

// Keys returns every key known to the counter in ascending order.
func (c *Counter) Keys() []string {
	c.mu.Lock()
	keys := make([]string, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	c.mu.Unlock()

	sort.Strings(keys)
	return keys
}

// Snapshot copies the counter into a plain map so callers can read it without
// holding on to internal state.
func (c *Counter) Snapshot() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make(map[string]int, len(c.m))
	for k, e := range c.m {
		out[k] = e.n
	}
	return out
}

// Reset drops every key, returning the counter to its initial state.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m = make(map[string]*entry)
}
