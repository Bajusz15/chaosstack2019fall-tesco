package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var postgres *sql.DB

func Init() error {
	var err error

	//connStr := "user=postgres dbname=postgres sslmode=verify-full"
	user := os.Getenv("PGUSER")
	database := os.Getenv("PGDATABASE")
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	connStr := "user=" + user + " dbname=" + database + " host=" + host + " port=" + port + " sslmode=disable"
	//connStr := "postgres://postgres:password@localhost/postgres?sslmode=disable" //os.Getenv("DB_CONNECTION")
	postgres, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = postgres.Ping(); err != nil {
		return err
	}
	fmt.Println("You are connected to your database")
	return nil
}
