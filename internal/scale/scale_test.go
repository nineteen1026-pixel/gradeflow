package scale

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		raw, max, want float64
	}{
		{45, 50, 0.9},
		{0, 50, 0},
		{-5, 50, 0},
		{60, 50, 1},
	}
	for _, c := range cases {
		got, err := Normalize(c.raw, c.max)
		if err != nil {
			t.Fatalf("Normalize(%v, %v) error = %v", c.raw, c.max, err)
		}
		if math.Abs(got-c.want) > 1e-9 {
			t.Fatalf("Normalize(%v, %v) = %v, want %v", c.raw, c.max, got, c.want)
		}
	}
}

func TestNormalizeErrors(t *testing.T) {
	if _, err := Normalize(10, 0); err == nil {
		t.Fatal("Normalize with zero max: want error, got nil")
	}
	if _, err := Normalize(math.NaN(), 10); err == nil {
		t.Fatal("Normalize with NaN: want error, got nil")
	}
}

func TestLetter(t *testing.T) {
	bands := Default()
	cases := []struct {
		score float64
		want  string
	}{
		{1.0, "A"},
		{0.90, "A"},
		{0.899, "B"},
		{0.75, "C"},
		{0.61, "D"},
		{0.10, "F"},
		{0.0, "F"},
	}
	for _, c := range cases {
		if got := Letter(c.score, bands); got != c.want {
			t.Fatalf("Letter(%v) = %q, want %q", c.score, got, c.want)
		}
	}
}

func TestLetterIgnoresBandOrder(t *testing.T) {
	shuffled := []Band{
		{Min: 0.60, Letter: "D"},
		{Min: 0.90, Letter: "A"},
		{Min: 0.00, Letter: "F"},
		{Min: 0.80, Letter: "B"},
	}
	if got := Letter(0.85, shuffled); got != "B" {
		t.Fatalf("Letter(0.85) = %q, want %q", got, "B")
	}
	if got := Letter(0.5, nil); got != "" {
		t.Fatalf("Letter with no bands = %q, want empty", got)
	}
}

func TestLetterDoesNotReorderCallerBands(t *testing.T) {
	bands := []Band{
		{Min: 0.60, Letter: "D"},
		{Min: 0.90, Letter: "A"},
	}
	_ = Letter(0.95, bands)

	if bands[0].Letter != "D" || bands[1].Letter != "A" {
		t.Fatalf("Letter reordered caller bands: %+v", bands)
	}
}

func TestRoundAndPercent(t *testing.T) {
	if got := Round(0.12345, 2); math.Abs(got-0.12) > 1e-9 {
		t.Fatalf("Round(0.12345, 2) = %v, want 0.12", got)
	}
	if got := Round(1.5, -1); math.Abs(got-2) > 1e-9 {
		t.Fatalf("Round(1.5, -1) = %v, want 2", got)
	}
	if got, want := Percent(0.8567), "85.7%"; got != want {
		t.Fatalf("Percent(0.8567) = %q, want %q", got, want)
	}
}
