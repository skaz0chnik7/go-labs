package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Шифрование AES-GCM
func encryptAES(plaintext, key string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("длина ключа должна быть 16, 24 или 32 символа")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ошибка создания шифра: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка создания GCM: %v", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("ошибка генерации nonce: %v", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Расшифровка AES-GCM
func decryptAES(ciphertext, key string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("длина ключа должна быть 16, 24 или 32 символа")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ошибка создания шифра: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка создания GCM: %v", err)
	}

	enc, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования Base64: %v", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("недопустимая длина шифротекста")
	}

	nonce, ciphertextBytes := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка расшифровки: %v", err)
	}

	return string(plaintext), nil
}

func main() {
	var option int

	fmt.Println("=== AES Шифрование/Расшифровка ===")
	fmt.Println("1. Зашифровать строку")
	fmt.Println("2. Расшифровать строку")
	fmt.Print("Выберите опцию: ")
	fmt.Scanln(&option)

	switch option {
	case 1:
		var plaintext, key string
		fmt.Print("Введите строку для шифрования: ")
		fmt.Scanln(&plaintext)
		fmt.Print("Введите ключ (16, 24 или 32 символа): ")
		fmt.Scanln(&key)

		encryptedText, err := encryptAES(plaintext, key)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Зашифрованные данные: %s\n", encryptedText)

	case 2:
		var ciphertext, key string
		fmt.Print("Введите строку для расшифровки: ")
		fmt.Scanln(&ciphertext)
		fmt.Print("Введите ключ: ")
		fmt.Scanln(&key)

		decryptedText, err := decryptAES(ciphertext, key)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Расшифрованные данные: %s\n", decryptedText)

	default:
		fmt.Println("Неверная опция. Завершение работы.")
	}
}
