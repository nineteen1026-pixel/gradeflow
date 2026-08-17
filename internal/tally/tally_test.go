package tally

import (
	"reflect"
	"testing"
)

func TestBumpCreatesAndAccumulates(t *testing.T) {
	c := New()

	if got := c.Bump("hw-1", 5); got != 5 {
		t.Fatalf("Bump(hw-1, 5) = %d, want 5", got)
	}
	if got := c.Bump("hw-1", 3); got != 8 {
		t.Fatalf("Bump(hw-1, 3) = %d, want 8", got)
	}
	if got := c.Bump("hw-2", 2); got != 2 {
		t.Fatalf("Bump(hw-2, 2) = %d, want 2", got)
	}
	if got := c.Value("hw-1"); got != 8 {
		t.Fatalf("Value(hw-1) = %d, want 8", got)
	}
}

func TestValueUnknownKey(t *testing.T) {
	c := New()
	if got := c.Value("nope"); got != 0 {
		t.Fatalf("Value(nope) = %d, want 0", got)
	}
}

func TestTotalAndKeys(t *testing.T) {
	c := New()
	c.Bump("b", 4)
	c.Bump("a", 6)
	c.Bump("c", -1)

	if got := c.Total(); got != 9 {
		t.Fatalf("Total() = %d, want 9", got)
	}
	if got, want := c.Keys(), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys() = %v, want %v", got, want)
	}
}

func TestSnapshotIsDetached(t *testing.T) {
	c := New()
	c.Bump("hw-1", 7)

	snap := c.Snapshot()
	if got, want := snap, map[string]int{"hw-1": 7}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Snapshot() = %v, want %v", got, want)
	}

	snap["hw-1"] = 100
	if got := c.Value("hw-1"); got != 7 {
		t.Fatalf("Value(hw-1) after mutating snapshot = %d, want 7", got)
	}
}

func TestReset(t *testing.T) {
	c := New()
	c.Bump("hw-1", 3)
	c.Reset()

	if got := c.Total(); got != 0 {
		t.Fatalf("Total() after Reset = %d, want 0", got)
	}
	if got := len(c.Keys()); got != 0 {
		t.Fatalf("len(Keys()) after Reset = %d, want 0", got)
	}
}
