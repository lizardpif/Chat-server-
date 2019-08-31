package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/CossackPyra/pyraconv"

	"./models"
)

//проверка существования таблиц, если нет, то создать!
//в передаче id текстовое или числовое!

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

func CheckAddChat(chat models.Chat) bool {
	//сделать проверку на повтор юзеров
	if !IsRepeate(chat.Users) {
		return false
	} //проверка а есть ли такой юзер?
	if !models.DbIsExistInts(chat.Users, "chat_data.users", "id") {
		return false
	} //если слишком мало юзеров
	if len(chat.Users) < 2 {
		return false
	} //а если такой чат уже существует?
	if !models.DbIsExistStr(chat.Name, "chat_data.chats_id", "name") {
		return false
	}
	return true
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
		}
		if !CheckAddChat(chat) {
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

	var msg models.Message

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	} else {
		err = json.Unmarshal(body, &msg)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		//существует ли чат?
		//существует ли автор?
		if !models.DbIsExistChatAuthor(msg.Chat, msg.Author, "chat_data.chats", "id") {

			http.Error(w, http.StatusText(400), 400)
			return
		}

	}
	id := models.DbMessageAdd(msg)
	if id == 0 {
		http.Error(w, http.StatusText(500), 500)
		return
	} else {
		fmt.Fprintln(w, id)
	}

}

func GetListOfChats(w http.ResponseWriter, r *http.Request) {
	//получить список чатов юзера по времени создания последнего сообщения
	if !IsPost(r.Method, w) {
		return
	}
	var user_id int
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "err %q\n", err, err.Error())
	} else {
		var result map[string]interface{}

		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		//fmt.Fprintln(w, result["user"])
		user_id = int(pyraconv.ToInt64(result["user"]))
	}
	chats := models.DbListOfChats(user_id)
	if chats == nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	fmt.Fprintln(w, chats)
}
func GetListOfMessages(w http.ResponseWriter, r *http.Request) {
	//получить список сообщений в чате по времени создания
	if !IsPost(r.Method, w) {
		return
	}
	var user_id int
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "err %q\n", err, err.Error())
	} else {
		var result map[string]interface{}

		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		//fmt.Fprintln(w, result["user"])
		user_id = int(pyraconv.ToInt64(result["chat"]))
	}
	messages := models.DbListOfMessages(user_id)
	// if chats == nil {
	// 	http.Error(w, http.StatusText(400), 400)
	// 	return
	// }
	fmt.Fprintln(w, messages)

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
