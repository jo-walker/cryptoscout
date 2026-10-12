package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Source struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Target struct {
	TimeoutSeconds int      `json:"timeoutSeconds"`
	Source         []Source `json:"source"`
}

func loadTargets(path string) (Target, error) {
	var cfg Target

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, fmt.Errorf("file %s does not exist: %w", path, err)
		}
		return cfg, fmt.Errorf("error reading file %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("%s is not a valid JSON: %w", path, err)
	}

	if len(cfg.Source) == 0 {
		return cfg, fmt.Errorf("%s has no targets", path)
	}

	for i, source := range cfg.Source {
		if source.Port == 0 {
			cfg.Source[i].Port = 443 // default port
		}
	}

	return cfg, nil
}

func main() {
	name := flag.String("config", "targets.json", "Path to the JSON configuration file")
	flag.Parse()

	cfg, err := loadTargets(*name)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error loading targets: %v\n", err))
		os.Exit(1)
	}
	for i, source := range cfg.Source {
		fmt.Printf("Target %d: %s:%d\n", i+1, source.Host, source.Port)
	}
}
