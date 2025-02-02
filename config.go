package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _config *Config = nil

type Environment struct {
	User string
	Pass string
	Host string
	Port string
}
type Credentials Environment
type Config struct {
	Credentials Credentials
	GateConfig  GateConfig
}
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
		}
		defer file.Close()

		var config Config
		b, err := io.ReadAll(file)
		if err != nil {
		}
		err = json.Unmarshal(b, &config)
		if err != nil {
			fmt.Println(config)
		}
		_config = &config
	}
	return _config
}
