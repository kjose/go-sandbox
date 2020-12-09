package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB

func Init() {
	db, err := sql.Open("mysql", "root:rootroot@tcp(go-sandbox-instance-1.cjt9ge78unj0.eu-west-3.rds.amazonaws.com:3306)/go-sandbox?charset=utf8")
	DB = db
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Init database connexion")
	defer db.Close()
}
