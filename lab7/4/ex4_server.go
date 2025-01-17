package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Data struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)

	fmt.Println("HTTP server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Случайная цитата Достоевского (просто потому что)
		fmt.Fprintln(w, "Боль и страдания всегда обязательны для широкого сознания и глубокого сердца!")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		fmt.Println("Полученный JSON:", data.Message)
		fmt.Fprintln(w, "Данные получены")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
