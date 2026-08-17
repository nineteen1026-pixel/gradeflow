package review

import (
	"errors"
	"reflect"
	"testing"
)

func seeded(t *testing.T) *Queue {
	t.Helper()

	q := NewQueue()
	for _, it := range []*Item{
		{ID: "r-1", Submission: "s-1", Priority: 1},
		{ID: "r-2", Submission: "s-2", Priority: 9},
		{ID: "r-3", Submission: "s-3", Priority: 5},
	} {
		if err := q.Push(it); err != nil {
			t.Fatalf("Push(%s) error = %v", it.ID, err)
		}
	}
	return q
}

func TestPushOrdersByPriority(t *testing.T) {
	q := seeded(t)

	if got := q.Len(); got != 3 {
		t.Fatalf("Len() = %d, want 3", got)
	}
	if got, want := q.IDs(), []string{"r-2", "r-3", "r-1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("IDs() = %v, want %v", got, want)
	}
}

func TestHeadReturnsFrontID(t *testing.T) {
	q := seeded(t)

	id, err := q.Head()
	if err != nil {
		t.Fatalf("Head() error = %v", err)
	}
	if id != "r-2" {
		t.Fatalf("Head() = %q, want %q", id, "r-2")
	}
	if got := q.Len(); got != 3 {
		t.Fatalf("Len() after Head = %d, want 3", got)
	}
}

func TestFetchDoesNotRemove(t *testing.T) {
	q := seeded(t)

	it, err := q.Fetch()
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if it.Submission != "s-2" {
		t.Fatalf("Fetch().Submission = %q, want %q", it.Submission, "s-2")
	}
	if got := q.Len(); got != 3 {
		t.Fatalf("Len() after Fetch = %d, want 3", got)
	}
}

func TestPopDrainsInOrder(t *testing.T) {
	q := seeded(t)

	var got []string
	for q.Len() > 0 {
		it, err := q.Pop()
		if err != nil {
			t.Fatalf("Pop() error = %v", err)
		}
		got = append(got, it.ID)
	}
	if want := []string{"r-2", "r-3", "r-1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("pop order = %v, want %v", got, want)
	}
	if _, err := q.Pop(); !errors.Is(err, ErrEmpty) {
		t.Fatalf("Pop() on drained queue error = %v, want ErrEmpty", err)
	}
}

func TestPushRejectsBadItems(t *testing.T) {
	q := NewQueue()
	if err := q.Push(nil); err == nil {
		t.Fatal("Push(nil): want error, got nil")
	}
	if err := q.Push(&Item{}); err == nil {
		t.Fatal("Push(item without id): want error, got nil")
	}
}
