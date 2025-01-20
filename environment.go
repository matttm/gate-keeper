package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var _environment *Environment = nil

type Environment struct {
	user   string
	pass   string
	host   string
	port   string
	config *GateConfig
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

func GetEnvironment() *Environment {
	if _environment == nil {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USERNAME")
		pass := os.Getenv("DB_PASSWORD")
		file, err := os.Open("config.json")
		if err != nil {
			log.Fatal("Error opening log file")
		}
		defer file.Close()

		var config GateConfig
		b, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Fatal(err)
			fmt.Println(config)
		}
		_environment = &Environment{
			user,
			pass,
			host,
			port,
			&config,
		}
	}
	return _environment
}
