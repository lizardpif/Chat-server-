// messages.go
package models

type Message struct {
	Chat   int    `json:"chat"`
	Author int    `json:"author"`
	Text   string `json:"text"`
}

func DbMessageAdd(message Message) int {
	//создать таблицу, если нет
	res, _ := db.Exec("INSERT INTO chat_data.messages (chat_id, author_id, text) VALUES (?,?,?)", message.Chat, message.Author, message.Text)
	id, _ := res.LastInsertId()

	return int(id)
}
