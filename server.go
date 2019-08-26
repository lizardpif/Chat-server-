package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Id        int    `json:"id"`
	UserName  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

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

	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Add new user")
		var user User

		//читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		//проверяем на наличие ошибки
		if err != nil {
			fmt.Fprintf(w, "err %q\n", err, err.Error())
		} else {
			//если все нормально - пишем по указателю в структуру
			err = json.Unmarshal(body, &user)
			if err != nil {
				fmt.Println(w, "can't unmarshal: ", err.Error())
			}
		}
		//выводим полученные данные (можно делать с данными все, что угодно)
		fmt.Fprintln(w, "id:", user.Id, "UserName:", user.UserName, "created_at:", user.CreatedAt)

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
