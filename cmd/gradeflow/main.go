// Command gradeflow runs a short end-to-end demo of the grading pipeline:
// submissions come in, points are tallied, a rubric is applied, the result is
// queued for human review and finally paged through.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gradeflow/internal/audit"
	"gradeflow/internal/page"
	"gradeflow/internal/review"
	"gradeflow/internal/rubric"
	"gradeflow/internal/scale"
	"gradeflow/internal/student"
	"gradeflow/internal/tally"
)

type submission struct {
	id      string
	student string
	points  float64
	max     float64
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gradeflow:", err)
		os.Exit(1)
	}
}

func run() error {
	logger := log.New(os.Stdout, "", 0)

	roster := student.NewRoster()
	for _, s := range []student.Student{
		{ID: "s-1", Name: "Ada", Cohort: "fall"},
		{ID: "s-2", Name: "Grace", Cohort: "fall"},
		{ID: "s-3", Name: "Alan", Cohort: "spring"},
	} {
		if err := roster.Add(s); err != nil {
			return err
		}
	}
	logger.Printf("enrolled %d students: %v", roster.Len(), roster.Names())

	sink := &audit.MemorySink{}
	trail := audit.New(sink)
	trail.SetLogger(logger)

	points := tally.New()
	queue := review.NewQueue()

	submissions := []submission{
		{id: "sub-1", student: "s-1", points: 46, max: 50},
		{id: "sub-2", student: "s-2", points: 39, max: 50},
		{id: "sub-3", student: "s-3", points: 44, max: 50},
		{id: "sub-4", student: "s-1", points: 31, max: 50},
	}

	logger.Println("--- scoring ---")
	for i, sub := range submissions {
		who, ok := roster.Get(sub.student)
		if !ok {
			return fmt.Errorf("submission %s references unknown student %s", sub.id, sub.student)
		}

		norm, err := scale.Normalize(sub.points, sub.max)
		if err != nil {
			return err
		}
		points.Bump(sub.student, int(sub.points))

		report, err := rubric.Apply(criteriaFor(norm))
		if err != nil {
			return err
		}

		letter := scale.Letter(report.Percent(), scale.Default())
		logger.Printf("%-6s %-22s raw=%s rubric=%s grade=%s lines=%v",
			sub.id, who.Display(), scale.Percent(norm), scale.Percent(report.Percent()), letter, report.Names())

		if err := trail.Record(audit.Entry{
			Actor:  "grader-bot",
			Action: "score",
			Target: sub.id,
			At:     time.Date(2024, 3, 1, 9, i, 0, 0, time.UTC),
		}); err != nil {
			return err
		}

		if err := queue.Push(&review.Item{
			ID:         "rev-" + sub.id,
			Submission: sub.id,
			Priority:   priorityFor(letter),
		}); err != nil {
			return err
		}
	}

	logger.Println("--- totals ---")
	for _, id := range points.Keys() {
		who, _ := roster.Get(id)
		logger.Printf("%-22s %d points", who.Display(), points.Value(id))
	}
	logger.Printf("course total: %d points", points.Total())

	logger.Println("--- review queue ---")
	head, err := queue.Head()
	if err != nil {
		return err
	}
	logger.Printf("%d items waiting, next up: %s", queue.Len(), head)

	const perPage = 2
	items, info := page.Window(queue.IDs(), 0, perPage)
	logger.Printf("page 1/%d (%d of %d): %v (more: %t)",
		page.Count(info.Total, perPage), info.Count, info.Total, items, info.HasNext)

	logger.Println("--- audit trail ---")
	for _, e := range sink.Entries() {
		logger.Printf("%s %s %s by %s", e.At.Format(time.RFC3339), e.Action, e.Target, e.Actor)
	}
	logger.Printf("%d audit entries written", trail.Written())

	return nil
}

func criteriaFor(norm float64) []rubric.Criterion {
	return []rubric.Criterion{
		{Name: "style", Weight: 0.2, Score: clamp01(norm + 0.05)},
		{Name: "correctness", Weight: 0.5, Score: norm},
		{Name: "tests", Weight: 0.3, Score: clamp01(norm - 0.10)},
	}
}

func clamp01(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}

func priorityFor(letter string) int {
	switch letter {
	case "F":
		return 9
	case "D":
		return 7
	case "C":
		return 5
	case "B":
		return 3
	default:
		return 1
	}
}
