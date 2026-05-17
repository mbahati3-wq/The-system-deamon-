package monitor

import (
    "context"
    "testing"
    "time"
)

func TestMonitorStartStop(t *testing.T) {
    cfg := &Config{
        CPUInterval:     50 * time.Millisecond,
        MemoryInterval:  50 * time.Millisecond,
        DiskInterval:    100 * time.Millisecond,
        NetworkInterval: 50 * time.Millisecond,
    }

    m, err := New(cfg)
    if err != nil {
        t.Fatalf("failed to create monitor: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
    defer cancel()

    if err := m.Start(ctx); err != nil {
        t.Fatalf("failed to start monitor: %v", err)
    }

    // Give it a moment to run
    time.Sleep(150 * time.Millisecond)

    if !m.IsRunning() {
        t.Fatalf("monitor should be running")
    }

    // Wait for context to expire and goroutine to stop
    <-ctx.Done()

    // allow some time for the monitor to process ctx.Done()
    time.Sleep(50 * time.Millisecond)

    if m.IsRunning() {
        t.Fatalf("monitor should have stopped after context done")
    }
}
