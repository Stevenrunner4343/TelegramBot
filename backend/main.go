package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
)

var token string = "12354535345asdfgh12314154sdfdgdh"
var wg = sync.WaitGroup{}

type User struct {
	Id         int    `json:"Id"`
	First_name string `json:"First_name"`
	Username   string `json:"Username"`
	Text       string `json:"Text"`
	Date       int    `json:"Date"`
}

type JsonData struct {
	Ok     bool `json:"ok"`
	Result []struct {
		Update_id int `json:"update_id"`
		Message   struct {
			Message_id int `json:"message_id"`
			From       struct {
				Id         int    `json:"id"`
				Is_bot     bool   `json:"is_bot"`
				First_name string `json:"first_name"`
				Last_name  string `json:"last_name"`
				Username   string `json:"username"`
			} `json:"from"`
			Chat struct {
				Id         int    `json:"id"`
				First_name string `json:"first_name"`
				Last_name  string `json:"last_name"`
				Username   string `json:"username"`
			}
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

func getDbData(db *sql.DB) {

	rows, err := db.Query("select * from users")
	if err != nil {
		fmt.Println("Ошибка при получении данных из БД", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		rows.Scan(&u.Id, &u.First_name, &u.Username, &u.Date)
		dbUsers[u.Id] = u

	}

}

var dbUsers = make(map[int]User)
var tgUsers []User

func getTgData(db *sql.DB) {
	TGcheck, err := http.Get("https://api.telegram.org/bot" + token + "/getUpdates")
	if err != nil {
		fmt.Println("Ошибка при запросе", err)
		fmt.Print("Статус:", TGcheck.Status)
	}
	var d JsonData
	json.NewDecoder(TGcheck.Body).Decode(&d)

	tgUsers = nil

	for i := range d.Result {

		u := User{
			Id:         d.Result[i].Message.From.Id,
			First_name: d.Result[i].Message.From.First_name,
			Username:   d.Result[i].Message.From.Username,
			Date:       d.Result[i].Message.Date,
			Text:       d.Result[i].Message.Text,
		}
		tgUsers = append(tgUsers, u)

	}

	for _, v := range tgUsers {
		if v.Id != dbUsers[v.Id].Id {
			dbUsers[v.Id] = v
			_, err := db.Exec(("insert into users (id, first_name, username, date) values ($1, $2, $3, $4)"), v.Id, v.First_name, v.Username, v.Date)
			if err != nil {
				fmt.Println("ошибка при записи в БД", err)
			}
			fmt.Println("Записан в мапу и Базу Данных")

		} else {
			fmt.Println("User уже есть")
		}

	}

}

var mapOfUsers = make(map[int]User)
var arrayOfMaps []User

func sendHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		getTgData(db)

		arrayOfMaps = nil
		for _, v := range tgUsers {
			mapOfUsers[v.Id] = v
			arrayOfMaps = append(arrayOfMaps, mapOfUsers[v.Id])

		}

		jsonData, err := json.Marshal(arrayOfMaps)
		if err != nil {
			fmt.Println("ошибка при кодировании в Json", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func recievehandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			fmt.Println("Ошибка дикедирования", err)
		}
		fmt.Println(u, "Собщение полученное с ФРОНТА")
		id := strconv.Itoa(u.Id)

		http.Get("https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + id + "&text=" + u.Text)

	}

}
func main() {
	connection := "user=appuser password=apppass dbname=appdb sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		fmt.Print("ошибка при создании подключения", err)
	}

	fmt.Println("Go запустилось и работает!")
	getDbData(db)

	http.HandleFunc("/getUsers", sendHandler(db))
	http.HandleFunc("/getMessages", recievehandler())

	err = http.ListenAndServe(":8084", nil)
	if err != nil {
		fmt.Println("ошибка при старте сервера", err)

	}

}
