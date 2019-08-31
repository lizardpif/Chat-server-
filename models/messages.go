// messages.go
package models

type Message struct {
	Id     int
	Chat   int    `json:"chat"`
	Author int    `json:"author"`
	Text   string `json:"text"`
	Date   string
}

func DbMessageAdd(message Message) int {
	//создать таблицу, если нет
	res, err := db.Exec("INSERT INTO chat_data.messages (chat_id, author_id, text) VALUES (?,?,?)", message.Chat, message.Author, message.Text)
	if err != nil {

		return 0
	}
	id, _ := res.LastInsertId()

	return int(id)
}
