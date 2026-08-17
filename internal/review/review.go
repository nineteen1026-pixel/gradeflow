// Package review models the queue of submissions waiting for a human grader.
package review

import (
	"errors"
	"sort"
)

// ErrEmpty is returned by queue operations that need at least one item.
var ErrEmpty = errors.New("review: queue is empty")

// Item is a single submission waiting to be reviewed. Higher Priority values
// are handed to graders first.
type Item struct {
	ID         string
	Submission string
	Priority   int
}

// Queue is an in-memory priority queue of review items. It is not safe for
// concurrent use; callers own the synchronisation.
type Queue struct {
	items []*Item
}

// NewQueue returns an empty queue.
func NewQueue() *Queue {
	return &Queue{}
}

// Push adds an item to the queue and keeps the queue ordered by priority.
// Items with equal priority keep their insertion order.
func (q *Queue) Push(it *Item) error {
	if it == nil {
		return errors.New("review: cannot push a nil item")
	}
	if it.ID == "" {
		return errors.New("review: item needs an id")
	}

	q.items = append(q.items, it)
	sort.SliceStable(q.items, func(i, j int) bool {
		return q.items[i].Priority > q.items[j].Priority
	})
	return nil
}

// Len reports how many items are waiting.
func (q *Queue) Len() int {
	return len(q.items)
}

// Fetch returns the highest priority item without removing it. It reports no
// item when the queue is empty.
func (q *Queue) Fetch() (*Item, error) {
	if len(q.items) == 0 {
		return nil, nil
	}
	return q.items[0], nil
}

// Pop removes and returns the highest priority item.
func (q *Queue) Pop() (*Item, error) {
	if len(q.items) == 0 {
		return nil, ErrEmpty
	}

	it := q.items[0]
	q.items = append(q.items[:0], q.items[1:]...)
	return it, nil
}

// Head reports the id of the submission currently at the front of the queue.
func (q *Queue) Head() (string, error) {
	it, err := q.Fetch()
	if err != nil {
		return "", err
	}
	return it.ID, nil
}

// IDs lists the queued item ids in review order.
func (q *Queue) IDs() []string {
	out := make([]string, 0, len(q.items))
	for _, it := range q.items {
		out = append(out, it.ID)
	}
	return out
}
