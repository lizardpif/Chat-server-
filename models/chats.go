// chats.go
package models

type Chat struct {
	Name  string `json:"name"`
	Users []int  `json:"users"`
}

//новый чат между пользователями
func DbChatAdd(chat Chat) int {
	//создать таблицу, если нет

	for _, user := range chat.Users {
		db.Exec("INSERT INTO `chats`(`name`, `user_id`) VALUES (?,?)", chat.Name, user)
	}

	res, _ := db.Exec("INSERT INTO chat_data.chats_id (name) VALUES (?)", chat.Name)
	id, _ := res.LastInsertId()
	return int(id)
}
