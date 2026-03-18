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
	EnablePprof bool // Enable pprof server on :8080 for profiling
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

func resolveConfigPath() (string, error) {
	if filePath := os.Getenv("GATE_KEEPER_CONFIG"); filePath != "" {
		return filePath, nil
	}

	workingDir, err := os.Getwd()
	if err == nil {
		filePath := filepath.Join(workingDir, "config.json")
		if _, statErr := os.Stat(filePath); statErr == nil {
			return filePath, nil
		}
	}

	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(executablePath), "config.json"), nil
}

// GetConfig loads the configuration from config.json (singleton).
func GetConfig() *Config {
	if _config == nil {
		filePath, err := resolveConfigPath()
		if err != nil {
			fmt.Println("Error resolving config path:", err)
			panic(err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening config file:", err)
			panic(err)
		}
		defer file.Close()

		var config Config
		b, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			panic(err)
		}
		err = json.Unmarshal(b, &config)
		if err != nil {
			fmt.Println("Error parsing config file:", err)
			panic(err)
		}
		_config = &config
	}
	return _config
}
