package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ipv4"
)

func main() {
	http.HandleFunc("/", mainPage)
	port := ":313"
	println("Server Listen on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	user := User{"Alex", "Net"}
	js, _ := json.Marshal(user)
	w.Write(js)
}
