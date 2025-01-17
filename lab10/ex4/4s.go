package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const validToken = "sometoken"

// authorizeHandler выдаёт токен сессии.
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] Авторизация клиента с адреса: %s\n", time.Now().Format(time.RFC3339), r.RemoteAddr)
	w.Write([]byte(validToken))
}

// dataHandler проверяет токен из заголовка Authorization.
func dataHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token != validToken {
		log.Printf("[%s] Неавторизованный доступ. Токен: '%s'\n", time.Now().Format(time.RFC3339), token)
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	// Успешный доступ
	log.Printf("[%s] Доступ предоставлен клиенту: %s\n", time.Now().Format(time.RFC3339), r.RemoteAddr)
	w.Write([]byte("Доступ предоставлен. Защищённые данные: {\"data\":\"Совершенно секретно\"}\n"))
}

func main() {
	// Пути к сертификатам
	serverCertPath := "server.crt"
	serverKeyPath := "server.key"
	clientCertPath := "client.crt"

	// Загрузка сертификата сервера
	serverCert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки сертификата сервера: %v", err)
	}

	// Настройка пула доверенных сертификатов клиента
	clientCACert, err := os.ReadFile(clientCertPath)
	if err != nil {
		log.Fatalf("Ошибка чтения сертификата клиента: %v", err)
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	// TLS конфигурация
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    clientCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      ":8444",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/data", dataHandler)

	log.Println("Сервер запущен на порту 8444...")
	if err := server.ListenAndServeTLS(serverCertPath, serverKeyPath); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
