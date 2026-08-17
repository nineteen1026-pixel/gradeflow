package tally_test

import (
    "sync"
    "testing"

    "gradeflow/internal/tally"
)

func TestVerifyBumpUnderConcurrency(t *testing.T) {
    const workers = 8
    const perWorker = 500

    c := tally.New()

    var wg sync.WaitGroup
    wg.Add(workers)
    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < perWorker; j++ {
                c.Bump("hw-1", 1)
            }
        }()
    }
    wg.Wait()

    if got, want := c.Value("hw-1"), workers*perWorker; got != want {
        t.Fatalf("Value(hw-1) after %d concurrent Bump calls = %d, want %d", want, got, want)
    }
}

func TestVerifyBumpAndReadUnderConcurrency(t *testing.T) {
    const workers = 4
    const perWorker = 400

    c := tally.New()

    var wg sync.WaitGroup
    wg.Add(workers * 2)
    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < perWorker; j++ {
                c.Bump("hw-2", 2)
            }
        }()
        go func() {
            defer wg.Done()
            for j := 0; j < perWorker; j++ {
                _ = c.Value("hw-2")
                _ = c.Total()
            }
        }()
    }
    wg.Wait()

    if got, want := c.Value("hw-2"), workers*perWorker*2; got != want {
        t.Fatalf("Value(hw-2) after concurrent Bump/Value = %d, want %d", got, want)
    }
}
