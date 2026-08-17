package audit

import (
	"testing"
	"time"
)

func fixedClock(l *Log, at time.Time) {
	l.now = func() time.Time { return at }
}

func TestRecordWritesToSink(t *testing.T) {
	sink := &MemorySink{}
	l := New(sink)
	at := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	fixedClock(l, at)

	if err := l.Record(Entry{Actor: "grader-1", Action: "score", Target: "sub-9"}); err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	got := sink.Entries()
	if len(got) != 1 {
		t.Fatalf("len(Entries()) = %d, want 1", len(got))
	}
	if got[0].Actor != "grader-1" || got[0].Action != "score" || got[0].Target != "sub-9" {
		t.Fatalf("Entries()[0] = %+v, unexpected fields", got[0])
	}
	if !got[0].At.Equal(at) {
		t.Fatalf("Entries()[0].At = %v, want %v", got[0].At, at)
	}
	if l.Written() != 1 {
		t.Fatalf("Written() = %d, want 1", l.Written())
	}
}

func TestRecordKeepsExplicitTimestamp(t *testing.T) {
	sink := &MemorySink{}
	l := New(sink)
	fixedClock(l, time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC))

	at := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := l.Record(Entry{Action: "reopen", Target: "sub-1", At: at}); err != nil {
		t.Fatalf("Record() error = %v", err)
	}
	if got := sink.Entries()[0].At; !got.Equal(at) {
		t.Fatalf("At = %v, want %v", got, at)
	}
}

func TestRecordRejectsEmptyAction(t *testing.T) {
	l := New(&MemorySink{})
	if err := l.Record(Entry{Actor: "grader-1"}); err == nil {
		t.Fatal("Record() without action: want error, got nil")
	}
}

func TestRecordWithoutSink(t *testing.T) {
	l := &Log{}
	if err := l.Record(Entry{Action: "score"}); err == nil {
		t.Fatal("Record() without sink: want error, got nil")
	}
}

func TestRecordAll(t *testing.T) {
	sink := &MemorySink{}
	l := New(sink)

	err := l.RecordAll([]Entry{
		{Actor: "a", Action: "score", Target: "sub-1"},
		{Actor: "b", Action: "score", Target: "sub-2"},
	})
	if err != nil {
		t.Fatalf("RecordAll() error = %v", err)
	}
	if l.Written() != 2 {
		t.Fatalf("Written() = %d, want 2", l.Written())
	}
	if len(sink.Entries()) != 2 {
		t.Fatalf("len(Entries()) = %d, want 2", len(sink.Entries()))
	}
}

func TestMemorySinkEntriesIsCopy(t *testing.T) {
	sink := &MemorySink{}
	if err := sink.Write(Entry{Action: "score"}); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	got := sink.Entries()
	got[0].Action = "mutated"
	if sink.Entries()[0].Action != "score" {
		t.Fatal("Entries() returned a view of internal state")
	}
}
