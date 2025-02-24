package test

import (
    "bytes"
    "net/http"
    "strings"
    "testing"
)

func TestEndToEnd(t *testing.T) {
    // Simulate scraping metrics
    resp, err := http.Get("http://localhost:8080/metrics")
    if err != nil {
        t.Fatalf("Failed to scrape metrics: %v", err)
    }
    defer resp.Body.Close()

    // Read the response body
    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    metrics := buf.String()

    // Load expected metrics from a fixture file
    expectedMetrics, err := ioutil.ReadFile("test/fixtures/expected_metrics.txt")
    if err != nil {
        t.Fatalf("Failed to read expected metrics file: %v", err)
    }

    // Compare actual metrics with expected metrics
    if !strings.Contains(metrics, string(expectedMetrics)) {
        t.Errorf("Scraped metrics do not match expected metrics")
    }
}
