package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"./models"
)

func main() {

	fmt.Println("Server is listening...")
	models.DbOpen("root:@/chat_data")
	defer models.DbClose()

	http.HandleFunc("/users/add", AddUser)
	http.HandleFunc("/chats/add", AddChat)
	http.HandleFunc("/chats/get", GetListOfChats)
	http.HandleFunc("/messages/add", AddMessage)
	http.HandleFunc("/messages/get", GetListOfMessages)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "messages/get messages/add chats/get chats/add users/add")
	})

	log.Fatal(http.ListenAndServe("localhost:8181", nil))
}

func IsPost(post string, w http.ResponseWriter) bool {
	if post != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return false
	}
	return true
}

func AddUser(w http.ResponseWriter, r *http.Request) {

	if !IsPost(r.Method, w) {
		return
	}
	var user models.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "err %q\n", err, err.Error())
	} else {
		err = json.Unmarshal(body, &user)
		if user.UserName == "" || err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}
	id := models.DbUserAdd(user.UserName)
	if id == 0 {
		http.Error(w, http.StatusText(500), 500)
		return
	} else {
		fmt.Fprintln(w, id)
	}

}

func AddChat(w http.ResponseWriter, r *http.Request) {
	if !IsPost(r.Method, w) {
		return
	}
	//новый чат между пользователями users_id
	var chat models.Chat

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	} else {
		err = json.Unmarshal(body, &chat)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		} //сделать проверку на повтор юзеров
		if !IsRepeate(chat.Users) {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}

	id := models.DbChatAdd(chat)
	if id == 0 {
		http.Error(w, http.StatusText(500), 500)
		return
	} else {
		fmt.Fprintln(w, id)
	}

}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	//отправить сообшение от юзера
	if !IsPost(r.Method, w) {
		return
	}
}
func GetListOfChats(w http.ResponseWriter, r *http.Request) {
	//получить список чатов юзера по времени создания
	if !IsPost(r.Method, w) {
		return
	}
}
func GetListOfMessages(w http.ResponseWriter, r *http.Request) {
	//получить список сообщений в чате по времени создания
	if !IsPost(r.Method, w) {
		return
	}
}

func IsRepeate(users []int) bool {
	sort.Ints(users)
	for i, us := range users {
		if i > 0 && us == users[i-1] {
			return false
		}
	}
	return true
}
