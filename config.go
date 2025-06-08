package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _config *Config = nil

// Environment holds database connection details.
type Environment struct {
	User string
	Pass string
	Host string
	Port string
}

// Credentials is an alias for Environment.
type Credentials Environment

// Config holds the application's configuration.
type Config struct {
	Credentials Credentials
	GateConfig  GateConfig
}

// GateConfig holds configuration for gate table/fields.
type GateConfig struct {
	Dbname               string
	TableName            string
	GateNameKey          string
	GateYearKey          string
	GateOrderKey         string
	GateIsApplicableFlag string
	StartKey             string
	EndKey               string
}

// GetConfig loads the configuration from config.json (singleton).
func GetConfig() *Config {
	if _config == nil {
		// Get the directory where the executable is located.
		executablePath, err := os.Executable()
		if err != nil {
			fmt.Println("Error getting executable path:", err)
			panic(err)
		}
		executableDir := filepath.Dir(executablePath)
		// Construct the path to the JSON file relative to the executable.
		filePath := filepath.Join(executableDir, "config.json") // or "config/data.json" if in a subdirectory
		file, err := os.Open(filePath)
		if err != nil {
			// Handle error opening config file
		}
		defer file.Close()

		var config Config
		b, err := io.ReadAll(file)
		if err != nil {
			// Handle error reading config file
		}
		err = json.Unmarshal(b, &config)
		if err != nil {
			fmt.Println(config)
		}
		_config = &config
	}
	return _config
}
