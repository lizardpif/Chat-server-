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

func main() {

	fmt.Println("Для работы с чат-сервером убедитесь, что mysql запущена")
	fmt.Println("В начале работы сервер открывает базу данных chat_data (!убедитесь, что у вас есть база данных с названием chat_data!) и создает (если нету) таблицы: users, messages, chats, chats_id.")
	fmt.Println("В случае, если произошла ошибка с соединением к базе данных, будет выведена ошибка.")
	fmt.Printf("Для работы с чат-сервером используются следующие маршруты:\n /users/add - добавить нового пользователя\n")
	fmt.Printf(" /chats/add - создание нового чата\n messages/add - создание нового сообщения \n chats/get - получить список всех чатов конкретного пользователя, отсортированных по времени создания последнего сообщения")
	fmt.Printf(",по убыванию\n messages/get - получить список всех сообщений в конкретном чате, отсортированных по времени создания последнего сообщения, по возрастанию\n")
	fmt.Printf("При работе с POST-запросами следующие ошибки:\n 400 - в запросе синтаксическая ошибка ('uer' вместо 'user' или если пользователь с таким именем уже существует)\n 500 - ошибка в запросе к базе данных\n")
	fmt.Println("405 - не соотвествие метода запроса (GET вместо POST, например)")

	fmt.Printf("\n\nServer is listening...\n")

	models.DbOpen("root:@/chat_data")
	models.DbCreateTables()

	defer models.DbClose()

	http.HandleFunc("/users/add", AddUser)
	http.HandleFunc("/chats/add", AddChat)
	http.HandleFunc("/chats/get", GetListOfChats)
	http.HandleFunc("/messages/add", AddMessage)
	http.HandleFunc("/messages/get", GetListOfMessages)

	log.Fatal(http.ListenAndServe("localhost:9000", nil))
}

func CheckAddChat(chat models.Chat) bool {

	//повтор юзеров
	if !IsRepeate(chat.Users) { //sintax
		return false
	} //проверка а есть ли такой юзер?
	if !models.DbIsExistInts(chat.Users, "chat_data.users", "id") {
		return false //not found!!! 404
	} //если слишком мало юзеров
	if len(chat.Users) < 2 { //sintax
		return false
	} //а если такой чат уже существует?
	if !models.DbIsExistStr(chat.Name, "chat_data.chats_id", "name") {
		return false
	}
	return true
}

//несоответсвие метода запроса
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
		http.Error(w, http.StatusText(400), 400)
		return
	} else {
		err = json.Unmarshal(body, &user)
		if user.UserName == "" || err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}
	if !models.DbIsExistStr(user.UserName, "chat_data.users", "username") {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	id := models.DbUserAdd(user.UserName)
	if id == 0 {
		http.Error(w, http.StatusText(500), 500) //ошибка в запросе к бд
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
		var tmp_chat models.TmpChat
		var users []int

		//тут
		err = json.Unmarshal(body, &tmp_chat)
		if tmp_chat.Name == "" || err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		for _, tmp := range tmp_chat.Users {
			users = append(users, int(pyraconv.ToInt64(tmp)))
		}
		chat.Users = users
		chat.Name = tmp_chat.Name
		//fmt.Fprintln(w, chat)

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
		http.Error(w, http.StatusText(400), 400) //синтаксическая ошибка
		return
	} else {

		var tmp_message models.TmpMessage

		//тут
		err = json.Unmarshal(body, &tmp_message)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		msg.Author = int(pyraconv.ToInt64(tmp_message.Author))
		msg.Chat = int(pyraconv.ToInt64(tmp_message.Chat))
		msg.Text = tmp_message.Text

		//существует ли чат?
		//существует ли автор?
		if !models.DbIsExistChatAuthor(msg.Chat, msg.Author, "chat_data.chats", "id") {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}
	id := models.DbMessageAdd(msg)
	if id == 0 {
		http.Error(w, http.StatusText(500), 500) //ошибка в запросе, возможно не совпадает название таблицы
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
		http.Error(w, http.StatusText(400), 400)
		return
	} else {
		var result map[string]interface{}

		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		//а есть ли такой юзер?
		user_id = int(pyraconv.ToInt64(result["user"]))
		if models.DbIsExistStr(pyraconv.ToString(result["user"]), "chat_data.users", "id") {
			//fmt.Fprintln(w, "err")
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}
	er, chats := models.DbListOfChats(user_id)
	if er != true {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if len(chats) > 0 {
		for _, ch := range chats {
			fmt.Fprintln(w, ch.PrintStruct())
		}
	}

}

func GetListOfMessages(w http.ResponseWriter, r *http.Request) {
	//получить список сообщений в чате по времени создания
	if !IsPost(r.Method, w) {
		return
	}
	var user_id int
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	} else {
		var result map[string]interface{}

		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		user_id = int(pyraconv.ToInt64(result["chat"]))
		if models.DbIsExistStr(pyraconv.ToString(result["chat"]), "chat_data.chats", "id") {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}

	er, messages := models.DbListOfMessages(user_id)
	if er != true {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, msg := range messages {
		fmt.Fprintln(w, msg.PrintStruct())
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
