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

type application struct {
	flags  *Flags
	config *Config
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
	for _, target := range app.config.Targets {
		systemReport, err := fetchReport(fmt.Sprintf("%s:%d", target.Host, target.Port))
		if err != nil {
			log.Printf("Error fetching report from target %s: %v\n", target.Name, err)
			continue
		}

		log.Printf("Report from target %s:\n%+v\n", target.Name, systemReport)
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
