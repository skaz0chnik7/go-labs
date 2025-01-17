package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/message", messageHandler)
	mux.HandleFunc("/data", dataHandler)

	// Добавляем middleware для логирования
	loggedMux := loggingMiddleware(mux)

	fmt.Println("HTTP server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method: %s, URL: %s, Time: %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Мир сдвинулся с места!")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение данных из тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// Вывод данных в консоль
	fmt.Println("Received data:", string(body))

	// Ответ клиенту
	fmt.Fprintln(w, "Data received and logged")
}
