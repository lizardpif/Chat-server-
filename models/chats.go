// chats.go
package models

type Chat struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Users     string `json:"users"`
	CreatedAt string `json:"created_at"`
}

//новый чат между пользователями
func DbChatAdd(s string) int {

}
