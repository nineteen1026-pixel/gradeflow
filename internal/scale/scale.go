// Package scale converts raw points into normalised scores and letter grades.
package scale

import (
	"fmt"
	"math"
	"sort"
)

// Band maps everything at or above Min (expressed as a 0..1 share) to a letter.
type Band struct {
	Min    float64
	Letter string
}

// Default returns the grading bands used by the standard course template.
func Default() []Band {
	return []Band{
		{Min: 0.90, Letter: "A"},
		{Min: 0.80, Letter: "B"},
		{Min: 0.70, Letter: "C"},
		{Min: 0.60, Letter: "D"},
		{Min: 0.00, Letter: "F"},
	}
}

// Normalize expresses raw points as a share of max, clamped to [0, 1].
func Normalize(raw, max float64) (float64, error) {
	if max <= 0 {
		return 0, fmt.Errorf("scale: max points must be positive, got %v", max)
	}
	if math.IsNaN(raw) {
		return 0, fmt.Errorf("scale: raw points is not a number")
	}

	v := raw / max
	switch {
	case v < 0:
		return 0, nil
	case v > 1:
		return 1, nil
	default:
		return v, nil
	}
}

// Letter maps a normalised score onto the highest band it reaches. Bands are
// considered from the highest threshold down, so their input order does not
// matter.
func Letter(score float64, bands []Band) string {
	if len(bands) == 0 {
		return ""
	}

	sorted := make([]Band, len(bands))
	copy(sorted, bands)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Min > sorted[j].Min })

	for _, b := range sorted {
		if score >= b.Min {
			return b.Letter
		}
	}
	return sorted[len(sorted)-1].Letter
}

// Round rounds v to the given number of decimal places.
func Round(v float64, places int) float64 {
	if places < 0 {
		places = 0
	}

	f := math.Pow(10, float64(places))
	return math.Round(v*f) / f
}

// Percent renders a normalised score as a human readable percentage.
func Percent(score float64) string {
	return fmt.Sprintf("%.1f%%", Round(score*100, 1))
}
