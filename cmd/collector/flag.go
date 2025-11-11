package main

import "flag"

// Flags holds all command-line flag values
type Flags struct {
	ConfigPath string
}

// ParseFlags parses command-line flags and returns a Flags struct
func ParseFlags() *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.ConfigPath, "config", "/etc/ruoke/config.yaml", "Path to configuration file")
	flag.Parse()
	return flags
}
