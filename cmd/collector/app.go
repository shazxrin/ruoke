package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"shazxrin.github.io/ruoke/internal/systemreport"
)

const (
	day  = 24
	hour = 60
	kb   = 1024
	mb   = 1024 * kb
	gb   = 1024 * mb
)

type application struct {
	flags  *Flags
	config *Config

	notifier Notifier
}

func (app *application) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(app.config.Interval) * time.Second)
	defer ticker.Stop()

	app.fetchReportsFromTargets()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down collector...")
			return nil
		case <-ticker.C:
			log.Println("Collecting reports from targets...")
			app.fetchReportsFromTargets()
		}
	}
}

func (app *application) fetchReportsFromTargets() {
	var msg string
	for _, target := range app.config.Targets {
		systemReport, err := fetchReport(fmt.Sprintf("%s:%d", target.Host, target.Port))
		if err != nil {
			log.Printf("Error fetching report from target %s: %v\n", target.Name, err)
			msg = msg + fmt.Sprintf("%s\nStatus: Down\n\n", target.Name)
			continue
		}

		uptimeDuration := time.Duration(systemReport.Uptime) * time.Second

		msg = msg + fmt.Sprintf(
			"%s\nStatus: Up\nUptime: %d d %d h %d m\nLoad: %.2f (1m) %.2f (5m) %.2f (15m)\nMemory: %.2f GB (U) %.2f GB (F) %.2f GB (T)\n\n",
			target.Name,
			int(uptimeDuration.Hours())/day, int(uptimeDuration.Hours())%day, int(uptimeDuration.Minutes())%hour,
			systemReport.Load1,
			systemReport.Load5,
			systemReport.Load15,
			float64(systemReport.UsedMemory)/gb,
			float64(systemReport.FreeMemory)/gb,
			float64(systemReport.TotalMemory)/gb,
		)
	}

	err := app.notifier.Notify("Status Report", msg)
	if err != nil {
		return
	}
}

func fetchReport(host string) (*systemreport.SystemReport, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("http://" + host + "/report")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch report from %s: %w", host, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("monitor returned non-OK status for %s: %d, body: %s", host, resp.StatusCode, string(bodyBytes))
	}

	var report systemreport.SystemReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("failed to decode JSON report from %s: %w", host, err)
	}

	return &report, nil
}
