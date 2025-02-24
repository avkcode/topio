package prometheus

import (
    "net/http"
    "testing"
)

func TestPrometheusServer(t *testing.T) {
    go StartPrometheusService("localhost", 8080)

    // Wait for the server to start
    time.Sleep(1 * time.Second)

    // Check if the /metrics endpoint is accessible
    resp, err := http.Get("http://localhost:8080/metrics")
    if err != nil {
        t.Fatalf("Failed to access /metrics endpoint: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status code 200, got %d", resp.StatusCode)
    }
}
