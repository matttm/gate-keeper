package main

import (
	"encoding/json"
	"log"
	"os"
)

var _environment *Environment = nil

type Environment struct {
	user   string
	pass   string
	host   string
	port   string
	config *Config
}
type Config struct {
	dbname         string
	tableName      string
	primaryKeyName string
	startKey       string
	endKey         string
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
		var b []byte
		var config *Config
		_, err = file.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(b, config)
		if err != nil {
			log.Fatal(err)
		}
		_environment = &Environment{
			user,
			pass,
			host,
			port,
			config,
		}
	}
	return _environment
}
