// Package rubric evaluates a submission against a weighted list of criteria.
package rubric

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNoCriteria is returned when a rubric carries no criteria at all.
var ErrNoCriteria = errors.New("rubric: no criteria supplied")

// Criterion is a single scored line of a rubric. Score is expressed on a 0..1
// scale and Weight is the share of the final grade the line is worth.
type Criterion struct {
	Name   string
	Weight float64
	Score  float64
}

// Line is one row of a rendered rubric report.
type Line struct {
	Name   string
	Weight float64
	Score  float64
	Points float64
}

// Report is the outcome of applying a rubric to a submission.
type Report struct {
	Lines  []Line
	Total  float64
	Weight float64
}

// Validate checks that a criterion is well formed.
func (c Criterion) Validate() error {
	if c.Name == "" {
		return errors.New("rubric: criterion needs a name")
	}
	if c.Weight < 0 {
		return fmt.Errorf("rubric: criterion %q has negative weight", c.Name)
	}
	if c.Score < 0 || c.Score > 1 {
		return fmt.Errorf("rubric: criterion %q score %.2f out of range", c.Name, c.Score)
	}
	return nil
}

// Apply scores a submission against the supplied criteria. The returned report
// lists the criteria from heaviest to lightest so the biggest contributors read
// first, and Total is the weighted sum of every line.
func Apply(in []Criterion) (Report, error) {
	if len(in) == 0 {
		return Report{}, ErrNoCriteria
	}
	for _, c := range in {
		if err := c.Validate(); err != nil {
			return Report{}, err
		}
	}

	// Work on a copy so the caller's slice is left untouched — Apply only reads
	// its input; the reordered view lives in the returned report.
	ordered := make([]Criterion, len(in))
	copy(ordered, in)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Weight > ordered[j].Weight
	})

	rep := Report{Lines: make([]Line, 0, len(ordered))}
	for _, c := range ordered {
		points := c.Weight * c.Score
		rep.Lines = append(rep.Lines, Line{
			Name:   c.Name,
			Weight: c.Weight,
			Score:  c.Score,
			Points: points,
		})
		rep.Total += points
		rep.Weight += c.Weight
	}
	return rep, nil
}

// Percent expresses the report total as a share of the available weight.
func (r Report) Percent() float64 {
	if r.Weight == 0 {
		return 0
	}
	return r.Total / r.Weight
}

// Names lists the criterion names in report order.
func (r Report) Names() []string {
	out := make([]string, 0, len(r.Lines))
	for _, l := range r.Lines {
		out = append(out, l.Name)
	}
	return out
}
