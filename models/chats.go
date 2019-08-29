// chats.go
package models

import "fmt"

type Chat struct {
	//Id        string `json:"id"`
	Name  string `json:"name"`
	Users []int  `json:"users"`
	//CreatedAt string `json:"created_at"`
}

//новый чат между пользователями
func DbChatAdd(chat Chat) int {

	var err error
	for i, user := range chat.Users {

		_, err = db.Exec("INSERT INTO `chats`(`name`, `user_id`) VALUES (?,?)", chat.Name, user)
		if err != nil {
			fmt.Println("чтото пошло не так с инсертом, итерация", i, user)
			fmt.Println(err)
			//err = nil
		}
	}

	res, _ := db.Exec("INSERT INTO chat_data.chats_id (name) VALUES (?)", chat.Name)

	id, _ := res.LastInsertId()

	return int(id)
}
