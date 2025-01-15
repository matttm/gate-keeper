package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitializeDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	// TODO: ADD VALIDATION
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/SPOTICLI_DB", user, pass, host, port),
	)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	flog.Successf("Database was successfully connected to")

	DB = db
}

func CloseDB() error {
	return DB.Close()
}

func GetDatabase() *sql.DB {
	if DB == nil {
		flog.Fatalf("Error: database not initialized")
	}
	return DB
}

//  func ExecSql[T any](s string) (T, error) {
//  	return GetDatabase().Exec(s)
//  }
