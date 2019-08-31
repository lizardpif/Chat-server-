// chats.go
package models

import "fmt"

type Chat struct {
	Id        int
	Name      string `json:"name"`
	Users     []int  `json:"users"`
	CreatedAt string
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

func DbListOfMessages(chat_id int) []Message {
	str := fmt.Sprintf("SELECT * FROM chat_data.messages WHERE chat_id=(?)")
	rows, err := db.Query(str, chat_id)
	if err != nil {
		return nil
	}
	var msgs []Message

	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Id, &msg.Chat, &msg.Author, &msg.Text, &msg.Date)
		if err != nil {
			return nil
		}
		//fmt.Println("msg ", msg)
		msgs = append(msgs, msg)
	}
	return msgs
}

func DbListOfChats(user_id int) []Chat {
	str := fmt.Sprintf("SELECT DISTINCT chats_id.id, name, chats_id.created_at FROM (chats_id JOIN chats USING (id)) JOIN messages ON chats_id.id=messages.chat_id WHERE chats.user_id=40 ORDER BY messages.created_at DESC")
	//
	rows, err := db.Query(str)
	if err != nil {
		//log.Panic(err)
		return nil
	}
	var chat Chat
	var chats []Chat
	var users []int
	var user int
	for rows.Next() {
		err = rows.Scan(&chat.Id, &chat.Name, &chat.CreatedAt)
		if err != nil {
			//fmt.Println("not scan 1")
			//return false
			return nil
		}
		row, er := db.Query("SELECT user_id FROM chat_data.chats WHERE id=(?)", chat.Id)
		for row.Next() {
			er = row.Scan(&user)
			if er != nil {
				//fmt.Println("not scan 2")
				//return false
				return nil
			}
			users = append(users, user)
		}
		chat.Users = users
		chats = append(chats, chat)
	}
	return chats
}
