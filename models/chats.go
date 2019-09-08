// chats.go
package models

type Chat struct {
	Id        int
	Name      string `json:"name"`
	Users     []int  `json:"users"`
	CreatedAt string
}

type TmpChat struct {
	Name  string
	Users []string
}

//новый чат между пользователями
func DbChatAdd(chat Chat) int {

	res, err := db.Exec("INSERT INTO chat_data.chats_id (name) VALUES (?)", chat.Name)
	if err != nil {
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	for _, user := range chat.Users {
		db.Exec("INSERT INTO `chats`(`id`, `user_id`) VALUES (?,?)", int(id), user)
	}

	return int(id)
}
