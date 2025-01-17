package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"log"
	"strings"
)

// Универсальная функция для вычисления хэша
func hashData(data, algorithm string) (string, error) {
	var hasher hash.Hash

	switch strings.ToLower(algorithm) {
	case "md5":
		hasher = md5.New()
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	default:
		return "", fmt.Errorf("неизвестный алгоритм хэширования: %s", algorithm)
	}

	// Вычисляем хэш
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("ошибка при вычислении хэша: %v", err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Функция проверки хэша
func verifyHash(data, providedHash, algorithm string) bool {
	calculatedHash, err := hashData(data, algorithm)
	if err != nil {
		log.Printf("Ошибка при вычислении хэша: %v", err)
		return false
	}
	return calculatedHash == providedHash
}

func main() {
	var input, algorithm string

	// Простое текстовое меню
	fmt.Println("=== Программа для работы с хэшами ===")
	fmt.Println("1. Вычислить хэш строки")
	fmt.Println("2. Проверить целостность данных")
	fmt.Println("Выберите опцию (1 или 2):")

	var option int
	fmt.Scanln(&option)

	switch option {
	case 1:
		fmt.Println("Введите строку для хэширования:")
		fmt.Scanln(&input)
		fmt.Println("Выберите алгоритм (md5, sha256, sha512):")
		fmt.Scanln(&algorithm)

		hash, err := hashData(input, algorithm)
		if err != nil {
			log.Fatalf("Ошибка: %v", err)
		}
		fmt.Printf("Хэш (%s): %s\n", algorithm, hash)

	case 2:
		fmt.Println("Введите строку для проверки:")
		fmt.Scanln(&input)
		fmt.Println("Введите хэш для сравнения:")
		var providedHash string
		fmt.Scanln(&providedHash)
		fmt.Println("Выберите алгоритм (md5, sha256, sha512):")
		fmt.Scanln(&algorithm)

		if verifyHash(input, providedHash, algorithm) {
			fmt.Println("Хэш совпадает! Данные не изменены.")
		} else {
			fmt.Println("Хэш не совпадает! Возможное повреждение данных.")
		}

	default:
		fmt.Println("Неверный выбор. Попробуйте снова.")
	}
}
