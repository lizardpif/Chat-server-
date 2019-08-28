// bd.go
package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DbOpen(sourceName string) {
	//"root:@/chat_data"
	var err error

	db, err = sql.Open("mysql", sourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}