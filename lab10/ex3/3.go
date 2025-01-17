package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Генерация пары ключей
func generateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// Сохранение ключа в файл
func saveKeyToFile(key interface{}, filename string, keyType string) error {
	var pemBlock *pem.Block

	switch k := key.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}
	case *rsa.PublicKey:
		keyBytes, err := x509.MarshalPKIXPublicKey(k)
		if err != nil {
			return err
		}
		pemBlock = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyBytes,
		}
	default:
		return errors.New("неизвестный тип ключа")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, pemBlock)
}

// Загрузка ключа из файла
func loadKeyFromFile(filename string, keyType string) (interface{}, error) {
	keyData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("неверный формат PEM в файле %s", filename)
	}

	switch keyType {
	case "private":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "public":
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return publicKey.(*rsa.PublicKey), nil
	default:
		return nil, errors.New("неизвестный тип ключа")
	}
}

// Подписание сообщения
func signMessage(privateKey *rsa.PrivateKey, message string) ([]byte, error) {
	hash := sha256.Sum256([]byte(message))
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

// Проверка подписи
func verifySignature(publicKey *rsa.PublicKey, message string, signature []byte) error {
	hash := sha256.Sum256([]byte(message))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}

func sender() {
	privateKey, err := generateRSAKeyPair(2048)
	if err != nil {
		log.Fatal("Ошибка генерации ключей:", err)
	}

	if err := saveKeyToFile(privateKey, "private.pem", "private"); err != nil {
		log.Fatal("Ошибка сохранения закрытого ключа:", err)
	}
	if err := saveKeyToFile(&privateKey.PublicKey, "public.pem", "public"); err != nil {
		log.Fatal("Ошибка сохранения открытого ключа:", err)
	}

	message := "Это сообщение от отправителя."
	signature, err := signMessage(privateKey, message)
	if err != nil {
		log.Fatal("Ошибка создания подписи:", err)
	}

	if err := ioutil.WriteFile("message.txt", []byte(message), 0644); err != nil {
		log.Fatal("Ошибка сохранения сообщения:", err)
	}
	if err := ioutil.WriteFile("signature.sig", signature, 0644); err != nil {
		log.Fatal("Ошибка сохранения подписи:", err)
	}

	fmt.Println("Сообщение и подпись успешно отправлены.")
}

func receiver() {
	publicKey, err := loadKeyFromFile("public.pem", "public")
	if err != nil {
		log.Fatal("Ошибка загрузки открытого ключа:", err)
	}

	message, err := ioutil.ReadFile("message.txt")
	if err != nil {
		log.Fatal("Ошибка чтения сообщения:", err)
	}
	signature, err := ioutil.ReadFile("signature.sig")
	if err != nil {
		log.Fatal("Ошибка чтения подписи:", err)
	}

	err = verifySignature(publicKey.(*rsa.PublicKey), string(message), signature)
	if err != nil {
		fmt.Println("Подпись недействительна!")
	} else {
		fmt.Println("Сообщение подлинно, подпись подтверждена.")
	}
}

func main() {
	var role int
	fmt.Println("Выберите роль: 1 - Отправитель, 2 - Получатель")
	fmt.Scanln(&role)

	switch role {
	case 1:
		sender()
	case 2:
		receiver()
	default:
		fmt.Println("Неверная роль.")
	}
}

// 3.	Асимметричное шифрование и цифровая подпись:
// •	Создайте пару ключей (открытый и закрытый) и сохраните их в файл.
// •	Реализуйте программу, которая подписывает сообщение с помощью закрытого ключа и проверяет подпись с использованием открытого ключа.
// •	Продемонстрируйте пример передачи подписанных сообщений между двумя сторонами.
