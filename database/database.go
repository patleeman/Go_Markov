package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
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
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, stmt := range data {
		_, err = tx.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

// Generate script to initialize database table for variable markov order.
func InitDB(markov_order int) {
	stmt := `DROP TABLE IF EXISTS markov; CREATE TABLE markov (m_id INTEGER PRIMARY KEY, target TEXT, %s);`
	var variable_columns string
	var col_name string
	for i := 0; i < markov_order; i++{
		col_name = "targetminus" + strconv.Itoa(markov_order - i) + " TEXT"
		if i != markov_order - 1 {
			col_name += ", "
		}
		variable_columns += col_name
	}
	stmt = fmt.Sprintf(stmt, variable_columns)
	ExecuteStatement(stmt)
}

