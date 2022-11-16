package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	db, err := sql.Open("mysql", "0nehaww9m3pq6ypvwvm2:pscale_pw_CwAqwhqSl8WXPibuybSu9heTJLheaw7cCQMNSqj1ZQY@tcp(us-east.connect.psdb.cloud)/cv-review?tls=true")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Database Connected!")

	DB = db
}
