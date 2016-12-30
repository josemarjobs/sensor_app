package datamanager

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	dbUser = "sensorappuser"
	dbPass = "sensorapppass"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres",
		"postgres://sensorappuser:sensorapppass@localhost/sensorapp?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
}
