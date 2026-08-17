// Package audit records grading decisions so they can be replayed later.
package audit

import (
	"errors"
	"io"
	"log"
	"time"
)

// Entry is a single audited action.
type Entry struct {
	Actor  string
	Action string
	Target string
	At     time.Time
}

// Sink is the destination an audit log writes to.
type Sink interface {
	Write(Entry) error
}

// MemorySink keeps entries in memory. It is handy for tests and for the
// development server.
type MemorySink struct {
	entries []Entry
}

// Write appends an entry to the sink.
func (m *MemorySink) Write(e Entry) error {
	m.entries = append(m.entries, e)
	return nil
}

// Entries returns a copy of everything written so far.
func (m *MemorySink) Entries() []Entry {
	out := make([]Entry, len(m.entries))
	copy(out, m.entries)
	return out
}

// Log writes audit entries to a sink and keeps a little bookkeeping.
type Log struct {
	sink    Sink
	logger  *log.Logger
	now     func() time.Time
	written int
}

// New returns a Log backed by sink. Diagnostics are discarded until a logger is
// installed with SetLogger.
func New(sink Sink) *Log {
	return &Log{
		sink:   sink,
		logger: log.New(io.Discard, "audit: ", 0),
		now:    time.Now,
	}
}

// SetLogger installs the logger used for diagnostics.
func (l *Log) SetLogger(lg *log.Logger) {
	if lg != nil {
		l.logger = lg
	}
}

// Record persists a single audit entry. Entries without a timestamp are stamped
// with the current time.
func (l *Log) Record(e Entry) error {
	if l.sink == nil {
		return errors.New("audit: log has no sink")
	}
	if e.Action == "" {
		return errors.New("audit: entry needs an action")
	}
	if e.At.IsZero() {
		e.At = l.now()
	}

	if err := l.sink.Write(e); err != nil {
		l.logger.Printf("write %s/%s failed: %v", e.Action, e.Target, err)
		return nil
	}

	l.written++
	return nil
}

// RecordAll writes every entry in order, stopping at the first failure.
func (l *Log) RecordAll(entries []Entry) error {
	for _, e := range entries {
		if err := l.Record(e); err != nil {
			return err
		}
	}
	return nil
}

// Written reports how many entries reached the sink.
func (l *Log) Written() int {
	return l.written
}
