package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shazxrin.github.io/ruoke/internal/systemreport"
	"time"
)

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

func main() {
	systemReport, err := fetchReport("localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(systemReport)
}
