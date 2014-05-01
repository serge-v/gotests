package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {

	db, err := sql.Open("mysql", "dbuser:pass1234@/test?parseTime=true")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer db.Close()

	rows, err := db.Query("select * from person")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for rows.Next() {
		var name string
		var sex string
		var birth time.Time
		var death time.Time

		err = rows.Scan(&name, &sex, &birth, &death)
		fmt.Printf("%-20s %s %s - %s\n", name, sex, birth.Format("2006-01-02"), death.Format("2006-01-02"))
	}
}
