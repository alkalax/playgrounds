package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type User struct {
	id   int
	name string
	age  int
}

func userExists(name string, db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT Id, Name, Age FROM Users WHERE Name = ?", name)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.id, &user.name, &user.age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func createUser(name string, age int, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO Users (name, age) VALUES (?, ?)", name, age)
	return err
}

func main() {
	db, err := sql.Open("sqlite", "practice.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Users (
		Id INTEGER PRIMARY KEY,
		Name TEXT,
		Age INTEGER
	)
	`); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	users, err := userExists(name, db)
	if err != nil {
		panic(err)
	}

	if len(users) > 0 {
		for i, user := range users {
			if i > 0 {
				fmt.Println("========================")
			}
			fmt.Println("User found.")
			fmt.Println("Id:", user.id)
			fmt.Println("Name:", user.name)
			fmt.Println("Age:", user.age)
		}
	} else {
		fmt.Printf("User %s not found. Creating...\n", name)
		fmt.Print("Age? ")
		ageStr, _ := reader.ReadString('\n')
		age, err := strconv.Atoi(strings.TrimSpace(ageStr))
		if err != nil {
			panic(err)
		}
		if err = createUser(name, age, db); err != nil {
			panic(err)
		}

		fmt.Println("User created.")
	}
}
