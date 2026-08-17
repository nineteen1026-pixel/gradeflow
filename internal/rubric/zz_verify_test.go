package rubric_test

import (
    "reflect"
    "testing"

    "gradeflow/internal/rubric"
)

func names(in []rubric.Criterion) []string {
    out := make([]string, 0, len(in))
    for _, c := range in {
        out = append(out, c.Name)
    }
    return out
}

func TestVerifyApplyKeepsCallerSliceOrder(t *testing.T) {
    in := []rubric.Criterion{
        {Name: "style", Weight: 0.2, Score: 1.0},
        {Name: "correctness", Weight: 0.5, Score: 0.8},
        {Name: "tests", Weight: 0.3, Score: 0.5},
    }
    before := names(in)

    if _, err := rubric.Apply(in); err != nil {
        t.Fatalf("Apply() error = %v", err)
    }

    if after := names(in); !reflect.DeepEqual(before, after) {
        t.Fatalf("Apply reordered the caller's slice: before %v, after %v", before, after)
    }
}

func TestVerifyApplyLeavesCallerElementsInPlace(t *testing.T) {
    in := []rubric.Criterion{
        {Name: "style", Weight: 0.2, Score: 1.0},
        {Name: "correctness", Weight: 0.5, Score: 0.8},
    }

    if _, err := rubric.Apply(in); err != nil {
        t.Fatalf("Apply() error = %v", err)
    }

    if in[0].Name != "style" || in[1].Name != "correctness" {
        t.Fatalf("caller slice after Apply = %v, want [style correctness]", names(in))
    }
}
