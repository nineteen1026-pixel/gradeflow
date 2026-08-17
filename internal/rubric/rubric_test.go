package rubric

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

func almost(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestApplyOrdersByWeight(t *testing.T) {
	rep, err := Apply([]Criterion{
		{Name: "style", Weight: 0.2, Score: 1.0},
		{Name: "correctness", Weight: 0.5, Score: 0.8},
		{Name: "tests", Weight: 0.3, Score: 0.5},
	})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if got, want := rep.Names(), []string{"correctness", "tests", "style"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Names() = %v, want %v", got, want)
	}
	if !almost(rep.Total, 0.75) {
		t.Fatalf("Total = %v, want 0.75", rep.Total)
	}
	if !almost(rep.Weight, 1.0) {
		t.Fatalf("Weight = %v, want 1.0", rep.Weight)
	}
	if !almost(rep.Percent(), 0.75) {
		t.Fatalf("Percent() = %v, want 0.75", rep.Percent())
	}
}

func TestApplyEmpty(t *testing.T) {
	if _, err := Apply(nil); !errors.Is(err, ErrNoCriteria) {
		t.Fatalf("Apply(nil) error = %v, want ErrNoCriteria", err)
	}
}

func TestApplyRejectsBadCriterion(t *testing.T) {
	if _, err := Apply([]Criterion{{Name: "style", Weight: 0.5, Score: 1.4}}); err == nil {
		t.Fatal("Apply() with out of range score: want error, got nil")
	}
	if _, err := Apply([]Criterion{{Name: "", Weight: 0.5, Score: 0.5}}); err == nil {
		t.Fatal("Apply() with unnamed criterion: want error, got nil")
	}
}

func TestValidate(t *testing.T) {
	if err := (Criterion{Name: "ok", Weight: 1, Score: 0.5}).Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	if err := (Criterion{Name: "bad", Weight: -1, Score: 0.5}).Validate(); err == nil {
		t.Fatal("Validate() with negative weight: want error, got nil")
	}
}

func TestPercentZeroWeight(t *testing.T) {
	rep, err := Apply([]Criterion{{Name: "bonus", Weight: 0, Score: 1}})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if rep.Percent() != 0 {
		t.Fatalf("Percent() = %v, want 0", rep.Percent())
	}
}
