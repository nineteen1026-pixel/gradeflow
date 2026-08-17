package student

import (
	"reflect"
	"testing"
)

func seeded(t *testing.T) *Roster {
	t.Helper()

	r := NewRoster()
	for _, s := range []Student{
		{ID: "s-1", Name: "Ada", Cohort: "fall"},
		{ID: "s-2", Name: "Grace", Cohort: "spring"},
		{ID: "s-3", Name: "Alan", Cohort: "fall"},
	} {
		if err := r.Add(s); err != nil {
			t.Fatalf("Add(%s) error = %v", s.ID, err)
		}
	}
	return r
}

func TestAddAndGet(t *testing.T) {
	r := seeded(t)

	if got := r.Len(); got != 3 {
		t.Fatalf("Len() = %d, want 3", got)
	}
	s, ok := r.Get("s-2")
	if !ok {
		t.Fatal("Get(s-2) not found")
	}
	if s.Name != "Grace" {
		t.Fatalf("Get(s-2).Name = %q, want %q", s.Name, "Grace")
	}
	if _, ok := r.Get("missing"); ok {
		t.Fatal("Get(missing) reported found")
	}
}

func TestAddRejectsBadInput(t *testing.T) {
	r := NewRoster()
	if err := r.Add(Student{Name: "No Id"}); err == nil {
		t.Fatal("Add without id: want error, got nil")
	}
	if err := r.Add(Student{ID: "s-1", Name: "  "}); err == nil {
		t.Fatal("Add without name: want error, got nil")
	}
	if err := r.Add(Student{ID: "s-1", Name: "Ada"}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if err := r.Add(Student{ID: "s-1", Name: "Ada Again"}); err == nil {
		t.Fatal("Add duplicate: want error, got nil")
	}
}

func TestIDsIsCopy(t *testing.T) {
	r := seeded(t)

	ids := r.IDs()
	if want := []string{"s-1", "s-2", "s-3"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("IDs() = %v, want %v", ids, want)
	}

	ids[0] = "clobbered"
	if r.IDs()[0] != "s-1" {
		t.Fatal("IDs() returned a view of internal state")
	}
}

func TestCohortSortedByName(t *testing.T) {
	r := seeded(t)

	got := r.Cohort("fall")
	if len(got) != 2 {
		t.Fatalf("len(Cohort(fall)) = %d, want 2", len(got))
	}
	if got[0].Name != "Ada" || got[1].Name != "Alan" {
		t.Fatalf("Cohort(fall) = %v, want Ada then Alan", []string{got[0].Name, got[1].Name})
	}
	if got := r.Cohort("summer"); got != nil {
		t.Fatalf("Cohort(summer) = %v, want nil", got)
	}
}

func TestNamesAndDisplay(t *testing.T) {
	r := seeded(t)

	if got, want := r.Names(), []string{"Ada", "Grace", "Alan"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Names() = %v, want %v", got, want)
	}
	if got, want := (Student{ID: "s-1", Name: "Ada", Cohort: "fall"}).Display(), "Ada (s-1, fall)"; got != want {
		t.Fatalf("Display() = %q, want %q", got, want)
	}
	if got, want := (Student{ID: "s-1", Name: "Ada"}).Display(), "Ada (s-1)"; got != want {
		t.Fatalf("Display() = %q, want %q", got, want)
	}
}
