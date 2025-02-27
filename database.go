package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitializeDatabase(user, pass, host, port, dbname string) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
	fmt.Printf("Connecting to %s...\n", url)
	db, err := sql.Open(
		"mysql",
		url,
	)
	if err != nil {
		fmt.Println("Error while connecting to database")
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		fmt.Println("Error while pinging database")
		panic(err)
	}
	log.Println("Database was successfully connected to")
	DB = db
}

func CloseDB() error {
	return DB.Close()
}

func GetDatabase() *sql.DB {
	if DB == nil {
		log.Fatalf("Error: database not initialized")
	}
	return DB
}

func ExecSql(s string) bool {
	_, err := GetDatabase().Exec(s)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
