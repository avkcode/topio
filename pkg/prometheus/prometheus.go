package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Metrics definitions
	processRchar = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_rchar_bytes",
			Help: "Number of bytes read from storage by the process.",
		},
		[]string{"pid", "label"},
	)
	processWchar = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_wchar_bytes",
			Help: "Number of bytes written to storage by the process.",
		},
		[]string{"pid", "label"},
	)
	processReadBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_read_bytes",
			Help: "Number of bytes actually read from storage by the process.",
		},
		[]string{"pid", "label"},
	)
	processWriteBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_write_bytes",
			Help: "Number of bytes actually written to storage by the process.",
		},
		[]string{"pid", "label"},
	)
	processCancelledWriteBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cancelled_write_bytes",
			Help: "Number of bytes cancelled during write operations by the process.",
		},
		[]string{"pid", "label"},
	)

	// Mutex to protect concurrent access to metrics
	metricsMutex sync.Mutex
)

// init registers the metrics with Prometheus
func init() {
	prometheus.MustRegister(processRchar)
	prometheus.MustRegister(processWchar)
	prometheus.MustRegister(processReadBytes)
	prometheus.MustRegister(processWriteBytes)
	prometheus.MustRegister(processCancelledWriteBytes)
}

// StartPrometheusService starts the Prometheus metrics server
func StartPrometheusService(bind string, port int) {
	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf("%s:%d", bind, port)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting Prometheus server: %v\n", err)
			os.Exit(1)
		}
	}()
	fmt.Printf("Prometheus metrics server started at http://%s:%d/metrics\n", bind, port)
}

// UpdatePrometheusMetrics updates the Prometheus metrics with process information
func UpdatePrometheusMetrics(info ProcessInfo) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	processRchar.WithLabelValues(info.PID, info.Label).Set(float64(info.Rchar))
	processWchar.WithLabelValues(info.PID, info.Label).Set(float64(info.Wchar))
	processReadBytes.WithLabelValues(info.PID, info.Label).Set(float64(info.ReadBytes))
	processWriteBytes.WithLabelValues(info.PID, info.Label).Set(float64(info.WriteBytes))
	processCancelledWriteBytes.WithLabelValues(info.PID, info.Label).Set(float64(info.CancelledWriteBytes))
}
