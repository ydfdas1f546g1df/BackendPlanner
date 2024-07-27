package dbt

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(dbuser, dbpassword, dbname, dbhost string) *sql.DB {
	var DB *sql.DB
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbuser, dbpassword, dbname, dbhost)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
	return DB
}

func CloseDB(DB *sql.DB) {
	DB.Close()
	fmt.Println("Successfully closed the database!")
}

func QueryDB(DB *sql.DB, query string) *sql.Rows {
	rows, err := DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func ExecDB(DB *sql.DB, query string) {
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
