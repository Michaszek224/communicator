package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type MessageStruct struct {
	Id      int
	Message string
	User    string
}

var templ *template.Template
var messages map[int]MessageStruct
var messagesMutex sync.RWMutex
var messageId int

func init() {
	messages = make(map[int]MessageStruct)
	messageId = 0
	var err error
	templ, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error passing templates: %v", err)
	}
}

func main() {
	http.HandleFunc("/", indexFunc)
	http.HandleFunc("/user", userFunc)
	fmt.Println("Server running on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexFunc(w http.ResponseWriter, r *http.Request) {
	err := templ.ExecuteTemplate(w, "index.html", messages)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func userFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			log.Printf("Error parsing form: %v", err)
			return
		}
		message := r.FormValue("message")
		user := r.FormValue("user")

		if message == "" || user == "" {
			http.Redirect(w, r, "/user", http.StatusSeeOther)
			return
		}

		messageLocalStruct := MessageStruct{
			Message: message,
			User:    user,
		}
		addMessage(messageLocalStruct)
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}
	err := templ.ExecuteTemplate(w, "user.html", nil)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}

}

func addMessage(msg MessageStruct) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()
	messageId++
	msg.Id = messageId
	messages[messageId] = msg
	fmt.Printf("Sending message:%v from: %v\n", msg.Message, msg.User)
	fmt.Printf("Current number of messages: %v\n", messageId)
}
