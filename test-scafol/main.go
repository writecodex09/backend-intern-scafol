package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Discussion struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desription"`
	Category    string `json:"category"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open(mysql.Open("root:@/scafol_database"), &gorm.Config{})

	if err != nil {
		log.Println("Connection Failed", err)
	} else {
		log.Println("Connection Established")
	}
	db.AutoMigrate(&Discussion{})

	handleRequests()
}
func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:8080")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/api/discussions", createDiscussion).Methods("POST")
	myRouter.HandleFunc("/api/discussions", getDiscussions).Methods("GET")
	myRouter.HandleFunc("/api/discussions/{id}", getDiscussion).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func createDiscussion(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var discussion Discussion
	json.Unmarshal(payloads, &discussion)

	db.Create(&discussion)

	res := Result{Code: 200, Data: discussion, Message: "Success create discussion"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getDiscussions(w http.ResponseWriter, r *http.Request) {
	discussions := []Discussion{}

	db.Find(&discussions)

	res := Result{Code: 200, Data: discussions, Message: "Success get all discussion"}

	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getDiscussion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	discussionID := vars["id"]

	var discussion Discussion

	db.First(&discussion, discussionID)

	res := Result{Code: 200, Data: discussion, Message: "Success get discussion by ID"}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
