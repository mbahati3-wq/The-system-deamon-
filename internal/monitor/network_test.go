package monitor

import (
    "testing"
)

func TestGetNetworkStats(t *testing.T) {
    stats, err := GetNetworkStats(false)
    if err != nil {
        t.Fatalf("GetNetworkStats returned error: %v", err)
    }

    if len(stats) == 0 {
        t.Fatalf("expected at least one network stat, got 0")
    }
}

func TestGetDetailedNetworkStats(t *testing.T) {
    data, err := GetDetailedNetworkStats()
    if err != nil {
        t.Fatalf("GetDetailedNetworkStats returned error: %v", err)
    }

    if data == nil {
        t.Fatalf("expected non-nil detailed network stats")
    }
}
