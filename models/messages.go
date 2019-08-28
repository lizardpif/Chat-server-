// messages.go
package models

type Message struct {
	Id        string `json:"id"`
	Chat      string `json:"chat"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}
