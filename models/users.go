// users.go
package models

type User struct {
	UserName string `json:"username"`
}

func DbUserAdd(u string) int {

	res, err := db.Exec("INSERT INTO chat_data.users (username) VALUES (?)", u)
	if err != nil {
		return 0
	}

	id, _ := res.LastInsertId()

	return int(id)

}
