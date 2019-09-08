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

func DbCreateTables() {

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, username VARCHAR(30), created_at DATETIME default CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal("users!")
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS chats (id INT NOT NULL, user_id INT NOT NULL)")
	if err != nil {
		log.Fatal("chats!")
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS chats_id (id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, name VARCHAR(30), created_at DATETIME default CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal("chats_id!")
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,chat_id INT NOT NULL, author_id INT NOT NULL, text TEXT, created_at DATETIME default CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal("messages!")
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

func DbIsExistInts(users []int, tablename string, param string) bool {

	str := fmt.Sprintf("SELECT * FROM %s WHERE %s=(?)", tablename, param)

	for _, us := range users {
		row, err := db.Query(str, us)
		defer row.Close()
		if err != nil {
			return false
		}
		if row.Next() {
			continue
		} else {
			return false
		}
	}
	return true
}

func DbIsExistStr(chatname string, tablename string, param string) bool {

	str := fmt.Sprintf("SELECT * FROM %s WHERE %s=(?)", tablename, param)
	row, err := db.Query(str, chatname)
	if err != nil {
		return false
	}
	defer row.Close()

	if row.Next() {
		return false
	} else {
		return true
	}

}

func DbIsExistChatAuthor(chat_id int, author_id int, tablename string, param string) bool {
	//если существует такой Id чата, проверить есть ли в нем id юзера

	str := fmt.Sprintf("SELECT * FROM %s WHERE %s=(?)", tablename, param)
	row, err := db.Query(str, chat_id)
	if err != nil {
		return false
	}
	defer row.Close()

	var scan_id int
	var ch_id int

	for row.Next() {

		err = row.Scan(&ch_id, &scan_id)
		if err != nil {

			return false
		}
		if author_id == scan_id {
			return true
		}
	}
	return false
}

func DbListOfMessages(chat_id int) (bool, []Message) {

	str := fmt.Sprintf("SELECT * FROM chat_data.messages WHERE chat_id=(?)")
	rows, err := db.Query(str, chat_id)

	if err != nil {
		return false, nil
	}
	defer rows.Close()

	var msgs []Message

	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Id, &msg.Chat, &msg.Author, &msg.Text, &msg.Date)
		if err != nil {
			return false, nil
		}
		msgs = append(msgs, msg)
	}
	return true, msgs
}

func DbListOfChats(user_id int) (bool, []Chat) {

	str := fmt.Sprintf("SELECT chats_id.id, name, chats_id.created_at FROM (chats_id JOIN chats USING (id)) JOIN messages ON chats_id.id=messages.chat_id WHERE chats.user_id=(?) ORDER BY messages.created_at DESC")
	rows, err := db.Query(str, user_id)
	if err != nil {
		return false, nil
	}
	defer rows.Close()

	var chat Chat
	var chats []Chat

	var user int
	var rep bool //переменная для проверки повторов

	for rows.Next() {
		err = rows.Scan(&chat.Id, &chat.Name, &chat.CreatedAt)
		if err != nil {
			return false, nil
		}
		var users []int
		row, er := db.Query("SELECT user_id FROM chat_data.chats WHERE id=(?)", chat.Id)
		if er != nil {
			return false, nil
		}
		defer row.Close()
		for row.Next() {
			er = row.Scan(&user)
			if er != nil {
				return false, nil
			}
			users = append(users, user)
		}
		chat.Users = users

		rep = true
		for _, ch := range chats { // проверка повторов

			if ch.Id == chat.Id {
				rep = false
			}
		}
		if rep {
			chats = append(chats, chat)
		}
	}

	return true, chats
}
