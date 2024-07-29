package main

import (
	"fmt"
	"github.com/mbfuss/sortingFiles/config"
	"github.com/mbfuss/sortingFiles/server"
	"log"
	"net/http"
	"os"
)

func main() {
	// Загружаем переменные из .env файла
	err := config.LoadEnv("config/serverPort.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Читаем порт из переменной окружения
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "80" // Значение по умолчанию, если PORT не задан
	}
	// Регистрация обработчика для пути /fs
	http.HandleFunc("/fs", server.HandleFileRequest)
	fmt.Printf("Сервер запущен на порту %s\n", port)
	// Запускает сервера, который будет прослушивать порт 80
	// nil -- использовать глобальный маршрутизатор http.DefaultServeMux для обработки запросов.
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
