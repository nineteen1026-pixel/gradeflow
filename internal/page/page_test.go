package page

import (
	"reflect"
	"testing"
)

var sample = []string{"a", "b", "c", "d", "e", "f"}

func TestWindowFirstPage(t *testing.T) {
	got, info := Window(sample, 0, 2)

	if want := []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Window(0, 2) = %v, want %v", got, want)
	}
	if info.Total != 6 || info.Count != 2 || info.Offset != 0 || info.Limit != 2 {
		t.Fatalf("Window(0, 2) info = %+v", info)
	}
	if !info.HasNext {
		t.Fatal("Window(0, 2) HasNext = false, want true")
	}
}

func TestWindowMiddlePage(t *testing.T) {
	got, info := Window(sample, 2, 2)

	if want := []string{"c", "d"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Window(2, 2) = %v, want %v", got, want)
	}
	if !info.HasNext {
		t.Fatal("Window(2, 2) HasNext = false, want true")
	}
}

func TestWindowNonPositiveLimit(t *testing.T) {
	got, info := Window(sample, 0, 0)
	if got != nil {
		t.Fatalf("Window(0, 0) = %v, want nil", got)
	}
	if info.Count != 0 || info.HasNext {
		t.Fatalf("Window(0, 0) info = %+v", info)
	}
}

func TestWindowOffsetPastEnd(t *testing.T) {
	got, info := Window(sample, 10, 3)
	if got != nil {
		t.Fatalf("Window(10, 3) = %v, want nil", got)
	}
	if info.Total != 6 || info.Count != 0 {
		t.Fatalf("Window(10, 3) info = %+v", info)
	}
}

func TestWindowNegativeOffset(t *testing.T) {
	got, info := Window(sample, -4, 2)
	if want := []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Window(-4, 2) = %v, want %v", got, want)
	}
	if info.Offset != 0 {
		t.Fatalf("Window(-4, 2).Offset = %d, want 0", info.Offset)
	}
}

func TestCount(t *testing.T) {
	cases := []struct {
		total, limit, want int
	}{
		{6, 2, 3},
		{7, 2, 4},
		{0, 2, 0},
		{6, 0, 0},
		{1, 10, 1},
	}
	for _, c := range cases {
		if got := Count(c.total, c.limit); got != c.want {
			t.Fatalf("Count(%d, %d) = %d, want %d", c.total, c.limit, got, c.want)
		}
	}
}

func TestOffsets(t *testing.T) {
	if got, want := Offsets(6, 2), []int{0, 2, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Offsets(6, 2) = %v, want %v", got, want)
	}
	if got := Offsets(0, 2); got != nil {
		t.Fatalf("Offsets(0, 2) = %v, want nil", got)
	}
}

func TestClamp(t *testing.T) {
	cases := []struct {
		limit, max, want int
	}{
		{5, 10, 5},
		{50, 10, 10},
		{0, 10, 1},
		{-3, 10, 1},
		{7, 0, 1},
	}
	for _, c := range cases {
		if got := Clamp(c.limit, c.max); got != c.want {
			t.Fatalf("Clamp(%d, %d) = %d, want %d", c.limit, c.max, got, c.want)
		}
	}
}
