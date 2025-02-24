package metrics

import (
    "strings"
    "testing"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetrics(t *testing.T) {
    info := ProcessInfo{
        PID:                  "1234",
        Label:                "test-process",
        Rchar:                1024,
        Wchar:                512,
        ReadBytes:            2048,
        WriteBytes:           1024,
        CancelledWriteBytes:  256,
    }

    UpdatePrometheusMetrics(info)

    // Test process_rchar_bytes
    if err := testutil.CollectAndCompare(
        processRchar,
        strings.NewReader(`
# HELP process_rchar_bytes Number of bytes read from storage by the process.
# TYPE process_rchar_bytes gauge
process_rchar_bytes{pid="1234",label="test-process"} 1024
`),
        "process_rchar_bytes",
    ); err != nil {
        t.Errorf("processRchar metric did not match: %v", err)
    }

    // Add similar tests for other metrics...
}
