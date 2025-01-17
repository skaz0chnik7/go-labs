package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var sessionToken string

// authorize выполняет авторизацию пользователя.
func authorize(client *http.Client, baseURL string) error {
	resp, err := client.Post(baseURL+"/authorize", "application/json", nil)
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа сервера: %v", err)
	}

	sessionToken = string(body)
	fmt.Printf("Токен сессии получен: %s\n", sessionToken)
	return nil
}

// sendRequestWithToken отправляет запрос с токеном.
func sendRequestWithToken(client *http.Client, baseURL string, token string) error {
	req, err := http.NewRequest("GET", baseURL+"/data", nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа сервера: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", string(body))
	return nil
}

func main() {
	// Пути к сертификатам
	clientCertPath := "client.crt"
	clientKeyPath := "client.key"
	serverCertPath := "server.crt"
	baseURL := "https://localhost:8444"

	// Загрузка клиентского сертификата
	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки клиентского сертификата: %v", err)
	}

	// Настройка пула доверенных сертификатов сервера
	serverCACert, err := os.ReadFile(serverCertPath)
	if err != nil {
		log.Fatalf("Ошибка чтения сертификата сервера: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(serverCACert)

	// TLS конфигурация
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	// Авторизация
	fmt.Println("Попытка авторизации...")
	if err := authorize(client, baseURL); err != nil {
		log.Fatalf("Ошибка авторизации: %v", err)
	}

	// Успешный запрос
	fmt.Println("Попытка отправить корректный запрос...")
	if err := sendRequestWithToken(client, baseURL, sessionToken); err != nil {
		log.Printf("Ошибка при отправке корректного запроса: %v", err)
	}

	// Некорректный запрос (с неверным токеном)
	fmt.Println("Попытка отправить некорректный запрос...")
	if err := sendRequestWithToken(client, baseURL, "invalid_token"); err != nil {
		log.Printf("Ошибка при отправке некорректного запроса: %v", err)
	}
}
