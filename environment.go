package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func GetEnvironment() *Config {
	if _config == nil {
		file, err := os.Open("config.json")
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
