package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./models"
)

// type Message struct {
// 	Id        string `json:"id"`
// 	Chat      string `json:"chat"`
// 	Author    string `json:"author"`
// 	Text      string `json:"text"`
// 	CreatedAt string `json:"created_at"`
// }

// type Chat struct {
// 	Id        string `json:"id"`
// 	Name      string `json:"name"`
// 	Users     string `json:"users"`
// 	CreatedAt string `json:"created_at"`
// }

func main() {

	fmt.Println("Server is listening...")
	models.DbOpen("root:@/chat_data")

	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Add new user")
		var user models.User

		//читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		//проверяем на наличие ошибки
		if err != nil {
			fmt.Fprintf(w, "err %q\n", err, err.Error())
		} else {
			//если все нормально - пишем по указателю в структуру
			err = json.Unmarshal(body, &user)
			if err != nil {
				fmt.Fprintln(w, "can't unmarshal: ", err.Error())
			}
		}
		fmt.Println("username = ", user.UserName)
		user.Id = models.DbUserAdd(user.UserName, user.CreatedAt)

		fmt.Println("id:", user.Id)

	})
	http.HandleFunc("/chats/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Add new chat with users")
	})
	http.HandleFunc("/chats/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Get list of all chats from user")
	})
	http.HandleFunc("/messages/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Send message from user in chat")
	})
	http.HandleFunc("/messages/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Get list of messages in chat")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "messages/get messages/add chats/get chats/add users/add")
	})

	log.Fatal(http.ListenAndServe("localhost:8181", nil))
}
