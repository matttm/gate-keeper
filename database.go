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
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USERNAME")
	// pass := os.Getenv("DB_PASSWORD")
	// TODO: ADD VALIDATION
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname),
	)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

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

//  func ExecSql[T any](s string) (T, error) {
//  	return GetDatabase().Exec(s)
//  }
