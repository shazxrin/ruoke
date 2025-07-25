package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"log"
	"net/http"
	"os"
	"shazxrin.github.io/ruoke/internal/systemreport"
	"time"
)

func createSystemReport() (*systemreport.SystemReport, error) {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	// Get uptime
	uptime, err := host.Uptime()
	if err != nil {
		return nil, fmt.Errorf("failed to get uptime: %w", err)
	}

	// Get load averages
	avg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get load averages: %w", err)
	}

	// Get memory usage
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	report := &systemreport.SystemReport{
		Hostname:    hostname,
		Uptime:      uptime,
		Load1:       avg.Load1,
		Load5:       avg.Load5,
		Load15:      avg.Load15,
		TotalMemory: vMem.Total,
		UsedMemory:  vMem.Used,
		FreeMemory:  vMem.Free,
		ReportTime:  time.Now().Format(time.RFC3339),
	}

	return report, nil
}

func getMetricsHandler(w http.ResponseWriter, r *http.Request) {
	report, err := createSystemReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(report); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	fmt.Println("Starting monitor at port 8080.")

	serveMux := http.ServeMux{}
	serveMux.HandleFunc("GET /report", getMetricsHandler)

	err := http.ListenAndServe(":8080", &serveMux)
	if err != nil {
		log.Fatal(err)
	}
}
