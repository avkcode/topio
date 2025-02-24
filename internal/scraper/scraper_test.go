package scraper

import (
    "bytes"
    "net/http"
    "strings"
    "testing"
)

func TestScrapeMetrics(t *testing.T) {
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

    // Check if specific metrics are present
    expectedMetrics := []string{
        `process_rchar_bytes{pid="1234",label="test-process"} 1024`,
        `process_wchar_bytes{pid="1234",label="test-process"} 512`,
    }

    for _, expected := range expectedMetrics {
        if !strings.Contains(metrics, expected) {
            t.Errorf("Expected metric not found: %s", expected)
        }
    }
}
