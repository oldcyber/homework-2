package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "qwe123"
	DB_NAME     = "item"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disabled", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Table struct {
	category string `json:"category_id"`
	user     string `json:"user_id"`
	date     string `json:"date_created"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Table `json:"data"`
	Message string  `json:"message"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/item/", GetTable).Methods("GET")

	router.HandleFunc("/item/", CreateNote).Methods("POST")

	router.HandleFunc("/item/{category_id}", DeleteNote).Methods("DELETE")

	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting table")

	rows, err := db.Query("SELECT * FROM Item")

	checkErr(err)

	var datas []Table

	for rows.Next() {
		var Category string
		var User string
		var Date string

		err = rows.Scan(&Category, &User, &Date)

		checkErr(err)

		datas = append(datas, Table{category: Category, user: User, date: Date})

	}

	var response = JsonResponse{Type: "success", Data: datas}

	json.NewEncoder(w).Encode(response)

}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("category_id")
	user := r.FormValue("user_id")
	date := r.FormValue("date_created")

	var response = JsonResponse{}

	if category == "" || user == "" || date == "" {
		response = JsonResponse{Type: "error", Message: "You are missing data."}
	} else {
		db := setupDB()

		printMessage("Inserting date into table")

		fmt.Println("Inserting new row in category: " + category + " by user: " + user + " at: " + date)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO Item(category_id, user_id, date_created) VALUES($1, $2, $3) returning id;", category, user, date).Scan(&lastInsertID)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The note has been added"}
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	category := params["category_id"]

	var response = JsonResponse{}

	if category == "" {
		response = JsonResponse{Type: "error", Message: "You trying to delete data that not exsist"}
	} else {
		db := setupDB()

		printMessage("Deleting note from DB")

		_, err := db.Exec("DELETE FROM Item where category_id = $1", category)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The note was deleted"}
	}

	json.NewEncoder(w).Encode(response)
}
