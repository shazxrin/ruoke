package systemreport

type SystemReport struct {
	Hostname    string  `json:"hostname"`
	Uptime      uint64  `json:"uptime_seconds"`
	Load1       float64 `json:"load_average_1min"`
	Load5       float64 `json:"load_average_5min"`
	Load15      float64 `json:"load_average_15min"`
	TotalMemory uint64  `json:"total_memory_bytes"`
	UsedMemory  uint64  `json:"used_memory_bytes"`
	FreeMemory  uint64  `json:"free_memory_bytes"`
	ReportTime  string  `json:"report_time"`
}
