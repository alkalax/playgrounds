package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "practice.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Users (
		Id INTEGER PRIMARY KEY,
		Name TEXT,
		Age INTEGER
	)
	`)
	if err != nil {
		panic(err)
	}
}
