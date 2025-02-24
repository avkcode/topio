package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ProcessInfo holds the extracted information for a process
type ProcessInfo struct {
	Timestamp           string
	PID                 string
	PPID                string
	Label               string
	Rchar               int64
	Wchar               int64
	ReadBytes           int64
	WriteBytes          int64
	CancelledWriteBytes int64
}

// getProcInfo extracts process information from /proc/<pid>/{cmdline,stat,io}
func getProcInfo(pid string) (ProcessInfo, error) {
	info := ProcessInfo{}

	// Read cmdline to get the process label
	cmdlinePath := filepath.Join("/proc", pid, "cmdline")
	cmdlineData, err := ioutil.ReadFile(cmdlinePath)
	if err != nil {
		return info, err
	}
	info.Label = strings.ReplaceAll(string(cmdlineData), "\x00", " ")

	// Read stat to get PPID
	statPath := filepath.Join("/proc", pid, "stat")
	statData, err := ioutil.ReadFile(statPath)
	if err != nil {
		return info, err
	}
	statFields := strings.Fields(string(statData))
	if len(statFields) > 3 {
		info.PPID = statFields[3]
	}

	// Read io to get rchar, wchar, read_bytes, write_bytes, cancelled_write_bytes
	ioPath := filepath.Join("/proc", pid, "io")
	ioData, err := ioutil.ReadFile(ioPath)
	if err != nil {
		return info, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(ioData))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			intValue, _ := strconv.ParseInt(value, 10, 64)
			switch key {
			case "rchar":
				info.Rchar = intValue
			case "wchar":
				info.Wchar = intValue
			case "read_bytes":
				info.ReadBytes = intValue
			case "write_bytes":
				info.WriteBytes = intValue
			case "cancelled_write_bytes":
				info.CancelledWriteBytes = intValue
			}
		}
	}

	info.PID = pid
	info.Timestamp = time.Now().UTC().Format(time.RFC3339)
	return info, nil
}

// listProcesses returns a list of all process IDs in /proc
func listProcesses() ([]string, error) {
	var pids []string
	entries, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			pid := entry.Name()
			if _, err := strconv.Atoi(pid); err == nil {
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}

func main() {
	// Parse command-line flags
	interval := flag.Int("interval", 1, "Polling interval in seconds")
	outputFile := flag.String("output", "", "Output file (default: stdout)")
	filterLabel := flag.String("filter", "", "Filter processes by label (e.g., 'wal writer')")
	prometheusBind := flag.String("prometheus-bind", "0.0.0.0", "Bind address for Prometheus metrics server")
	prometheusPort := flag.Int("prometheus-port", 9090, "Port for Prometheus metrics server")
	flag.Parse()

	// Start Prometheus service if enabled
	if *prometheusBind != "" && *prometheusPort > 0 {
		StartPrometheusService(*prometheusBind, *prometheusPort)
	}

	// Create a CSV writer
	var writer *csv.Writer
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = csv.NewWriter(file)
	} else {
		writer = csv.NewWriter(os.Stdout)
	}
	defer writer.Flush()

	// Write the header
	header := []string{
		"timestamp", "pid", "ppid", "label", "rchar", "wchar", "read_bytes", "write_bytes", "cancelled_write_bytes",
	}
	writer.Write(header)

	// Map to store previous values for delta calculations
	prevStats := make(map[string]ProcessInfo)

	for {
		// List all processes
		pids, err := listProcesses()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing processes: %v\n", err)
			os.Exit(1)
		}

		// Iterate through each process and extract information
		for _, pid := range pids {
			info, err := getProcInfo(pid)
			if err != nil {
				continue // Skip processes that cannot be read
			}

			// Apply filter if specified
			if *filterLabel != "" && !strings.Contains(info.Label, *filterLabel) {
				continue
			}

			// Calculate deltas if previous stats exist
			if prev, exists := prevStats[pid]; exists {
				info.Rchar -= prev.Rchar
				info.Wchar -= prev.Wchar
				info.ReadBytes -= prev.ReadBytes
				info.WriteBytes -= prev.WriteBytes
				info.CancelledWriteBytes -= prev.CancelledWriteBytes
			}
			prevStats[pid] = info

			// Write the row
			row := []string{
				info.Timestamp,
				info.PID,
				info.PPID,
				info.Label,
				strconv.FormatInt(info.Rchar, 10),
				strconv.FormatInt(info.Wchar, 10),
				strconv.FormatInt(info.ReadBytes, 10),
				strconv.FormatInt(info.WriteBytes, 10),
				strconv.FormatInt(info.CancelledWriteBytes, 10),
			}
			writer.Write(row)

			// Update Prometheus metrics
			if *prometheusBind != "" && *prometheusPort > 0 {
				UpdatePrometheusMetrics(info)
			}
		}

		// Sleep for the specified interval
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
