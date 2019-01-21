package main

import (
	"database/sql"
	_ "github.com/lib/pq"
  "os"
  "fmt"
  "log"
)

var db *sql.DB

func ConnectDB() {
  pgUser := os.Getenv("PG_USER")
  pgDbName := os.Getenv("PG_DB_NAME")
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", pgUser, pgDbName)

  // Create connection pointer
  var err error
  db, err = sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }

  // Actually establish connection
  err = db.Ping()
  if err != nil {
    log.Print("Error connecting to database:")
    log.Fatal(err)
  }
}
