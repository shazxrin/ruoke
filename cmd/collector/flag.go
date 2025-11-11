package main

import "flag"

type Flags struct {
	ConfigPath string
}

func ParseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.ConfigPath, "config", "/etc/ruoke/config.yaml", "Path to configuration file")

	flag.Parse()
	return flags
}
