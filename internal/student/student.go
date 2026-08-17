// Package student holds the roster of learners a course tracks.
package student

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Student is one enrolled learner.
type Student struct {
	ID     string
	Name   string
	Cohort string
}

// Display renders a student for logs and CLI output.
func (s Student) Display() string {
	if s.Cohort == "" {
		return fmt.Sprintf("%s (%s)", s.Name, s.ID)
	}
	return fmt.Sprintf("%s (%s, %s)", s.Name, s.ID, s.Cohort)
}

// Roster is an ordered collection of students keyed by id.
type Roster struct {
	byID  map[string]Student
	order []string
}

// NewRoster returns an empty roster.
func NewRoster() *Roster {
	return &Roster{byID: make(map[string]Student)}
}

// Add enrols a student. Ids must be unique and non-empty.
func (r *Roster) Add(s Student) error {
	if s.ID == "" {
		return errors.New("student: id must not be empty")
	}
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("student: %q has no name", s.ID)
	}
	if _, dup := r.byID[s.ID]; dup {
		return fmt.Errorf("student: %q already enrolled", s.ID)
	}

	r.byID[s.ID] = s
	r.order = append(r.order, s.ID)
	return nil
}

// Get looks a student up by id.
func (r *Roster) Get(id string) (Student, bool) {
	s, ok := r.byID[id]
	return s, ok
}

// Len reports how many students are enrolled.
func (r *Roster) Len() int {
	return len(r.order)
}

// IDs returns the enrolled ids in enrolment order.
func (r *Roster) IDs() []string {
	out := make([]string, len(r.order))
	copy(out, r.order)
	return out
}

// Cohort returns every student in the named cohort, sorted by name.
func (r *Roster) Cohort(name string) []Student {
	var out []Student
	for _, id := range r.order {
		if s := r.byID[id]; s.Cohort == name {
			out = append(out, s)
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Names lists every student name in enrolment order.
func (r *Roster) Names() []string {
	out := make([]string, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.byID[id].Name)
	}
	return out
}
