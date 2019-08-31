// bd.go
package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DbOpen(sourceName string) {
	//создание таблицы, если нету
	var err error

	db, err = sql.Open("mysql", sourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

func DbClose() {
	db.Close()
}

type PrintOut interface {
	PrintStruct() string
}

func PrintStruct(p PrintOut) {
	p.PrintStruct()
}

func (m Message) PrintStruct() string {
	str := fmt.Sprintf("id: %d, chat_id: %d, author_id: %d, text: %s, created_at: %s", m.Id, m.Chat, m.Author, m.Text, m.Date)
	return str
}
func (m Chat) PrintStruct() string {
	str := fmt.Sprintf("id: %d, name: %s, users: %d, created_at: %s", m.Id, m.Name, m.Users, m.CreatedAt)
	return str
}
