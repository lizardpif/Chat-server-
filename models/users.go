// users.go
package models

type User struct {
	Id        int    `json:"id"`
	UserName  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

func DbUserAdd(u string) int {

	res, err := db.Exec("INSERT INTO chat_data.users (username) VALUES (?)", u)
	if err != nil {
		return 0
	}

	id, _ := res.LastInsertId()

	return int(id)

}
