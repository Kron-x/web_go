package config

import (
    "encoding/json"
    "log"
    "os"
    "path/filepath"
)


type Config struct {
    Port          string `json:"port"`
    ImagesDir     string `json:"images_dir"`
    LogFile       string `json:"log_file"`
    PortMetrics   string `json:"port_metrics"`
}

func LoadConfig() Config {
    configPath := filepath.Join("configs", "config.json")
    file, err := os.ReadFile(configPath)
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }
    var config Config
    err = json.Unmarshal(file, &config)
    if err != nil {
        log.Fatalf("Failed to parse config file: %v", err)
    }
    return config
}