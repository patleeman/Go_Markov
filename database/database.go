package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"fmt"
)

func check(err error){
	if err != nil {
		log.Fatal(err)
	}
}

// Execute a sql statement
func ExecuteStatement(statement string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite3")
	check(err)
	defer db.Close()

	sqlStmt := statement

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func ExecuteTransaction(data []string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite3")
	check(err)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, stmt := range data {
		fmt.Println(stmt)
		_, err = tx.Exec(stmt)
		if err != nil {
			log.Fatal(err, " Statement: ", stmt)
		}
	}
	tx.Commit()
}

func Query(statement string, markov_order int) []string {
	// Open db
	db, err := sql.Open("sqlite3", "./database/db.sqlite3")
	check(err)

	rows, err := db.Query(statement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var options []string
	for rows.Next() {
		var target string
		rows.Scan(&target)

		if strings.Contains(target, "'"){
			target = strings.Replace(target, "'", "/'", -1)
		}

		options = append(options, target)
	}
	return options
}