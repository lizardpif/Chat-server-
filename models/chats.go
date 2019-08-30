// chats.go
package models

import "fmt"

type Chat struct {
	Name  string `json:"name"`
	Users []int  `json:"users"`
}

//новый чат между пользователями
func DbChatAdd(chat Chat) int {
	//создать таблицу, если нет
	res, _ := db.Exec("INSERT INTO chat_data.chats_id (name) VALUES (?)", chat.Name)
	id, _ := res.LastInsertId()
	for _, user := range chat.Users {
		db.Exec("INSERT INTO `chats`(`id`, `user_id`) VALUES (?,?)", int(id), user)
	}

	return int(id)
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
	defer row.Close()
	if err != nil {
		return false
	}
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

		err = row.Scan(&scan_id, &ch_id)
		if err != nil {
			return false
		}
		if author_id == scan_id {
			return true
		}
	}
	return false
}
