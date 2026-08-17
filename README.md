# gradeflow

`gradeflow` is a small, dependency-free Go service library for moving coursework
through a grading pipeline: a submission arrives, points are tallied, a weighted
rubric turns those points into a report, the result is queued for a human
reviewer, every decision is written to an audit trail, and the results are
served back to the UI one page at a time.

It is deliberately plain Go — no frameworks, no code generation, no external
modules — so it can be embedded in whatever HTTP or queue layer a course
platform already runs.

## Layout

```
cmd/gradeflow      end-to-end demo of the pipeline
internal/audit     append-only trail of grading decisions
internal/page      fixed-size windows over a result set
internal/review    priority queue of submissions awaiting a grader
internal/rubric    weighted criteria -> scored report
internal/scale     raw points -> normalised scores and letter grades
internal/student   the course roster
internal/tally     concurrency-safe point totals
```

## Getting started

```sh
go build ./...
go test ./...
go run ./cmd/gradeflow
```

The demo command enrols three students, scores four submissions, prints the
rubric breakdown and the letter grade for each, pushes everything onto the
review queue, and then pages through the queue and the audit trail.

## Package tour

### `internal/tally`

`tally.Counter` accumulates integer point totals keyed by student or assignment
id. It is shared by every grading worker, so all of its operations are safe to
call from multiple goroutines.

```go
c := tally.New()
c.Bump("s-1", 46)
c.Bump("s-1", 31)
c.Value("s-1") // 77
```

### `internal/rubric`

A rubric is a list of weighted `Criterion` values, each scored on a 0..1 scale.
`Apply` validates them, orders the resulting report from heaviest criterion to
lightest, and reports the weighted total.

```go
report, err := rubric.Apply([]rubric.Criterion{
    {Name: "style", Weight: 0.2, Score: 1.0},
    {Name: "correctness", Weight: 0.5, Score: 0.8},
    {Name: "tests", Weight: 0.3, Score: 0.5},
})
report.Percent() // 0.75
```

### `internal/review`

`review.Queue` keeps submissions ordered by priority, highest first, with stable
ordering for ties. `Fetch` peeks, `Pop` removes, and `Head` reports the id of
whatever a grader would be handed next.

### `internal/audit`

`audit.Log` writes `Entry` values to a `Sink`. `MemorySink` ships with the
package for tests and local development; production deployments plug in their
own database or message-bus sink.

### `internal/page`

`page.Window` slices a result set into fixed-size windows and hands back an
`Info` describing the window (offset, limit, total, count, and whether more
results follow). `page.Count`, `page.Offsets` and `page.Clamp` cover the
arithmetic API handlers usually re-implement by hand.

### `internal/scale` and `internal/student`

Supporting packages: normalising raw points against a maximum, mapping a
normalised score onto configurable letter-grade bands, and holding the roster of
enrolled learners.

## Requirements

Go 1.22 or newer. There are no third-party dependencies.

## License

MIT.
