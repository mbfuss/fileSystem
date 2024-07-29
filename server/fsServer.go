package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mbfuss/sortingFiles/config"
	"github.com/mbfuss/sortingFiles/service"
)

// HandleFileRequest - функция, которая обрабатывает http запросы
// http.ResponseWriter — интерфейс, который предоставляет методы для формирования и отправки HTTP-ответов клиенту
// http.Request — это указатель на структуру http.Request в Go, которая представляет собой HTTP-запрос, отправленный клиентом на сервер
func HandleFileRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Получение параметров запроса
	root, sortOrder, err := service.GetRequestParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Получение списка файлов из каталога
	entries, err := os.ReadDir(root)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка чтения директории: %v", err), http.StatusInternalServerError)
		return
	}

	// Обработка файлов и директорий для получения размеров
	fileInfoWithSizes := service.ProcessFiles(root, entries)
	// Сортировка файлов и директорий по размеру
	service.SortFiles(fileInfoWithSizes, sortOrder)

	// Указываем, что вывод будет в формате json
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w) создает новый JSON-энкодер, который будет записывать данные прямо в http.ResponseWriter (w), то есть отправлять данные клиенту
	// .Encode(fileInfoWithSizes) преобразует fileInfoWithSizes (срез структур FileInfoWithSize) в JSON-формат и отправляет его в ответе
	err = json.NewEncoder(w).Encode(fileInfoWithSizes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка кодирования в JSON: %v", err), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	fmt.Printf("Обработка запроса заняла: %v\n", duration)
}

func ServerStart() {
	// Загружаем переменные из .env файла
	err := config.LoadEnv("config/serverPort.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Читаем порт из переменной окружения
	port := os.Getenv("SERVER_PORT")

	// Создание нового сервера
	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	// Регистрация обработчика для пути /fs
	http.HandleFunc("/fs", HandleFileRequest)

	// Канал для получения системных сигналов
	stop := make(chan os.Signal, 1)
	// Функция signal.Notify регистрирует канал stop для получения уведомлений о сигнале os.Interrupt (Ctrl+C) и os.Kill (сигнал принудительного завершения процесса)
	// При получении одного из этих сигналов, он будет отправлен в канал stop
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Запуск сервера в отдельной горутине
	go func() {
		fmt.Printf("Сервер запущен на порту %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Сервер остановлен с ошибкой: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	fmt.Println("Получен сигнал завершения, остановка сервера...")

	//  Если сервер не завершится в течение 5 секунд, контекст будет отменён, и остановка сервера прервется
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Используется для отложенного вызова функции cancel. Это значит, что функция cancel будет вызвана автоматически,
	// когда выполнение функции ServerStart завершится, даже если произойдёт ошибка. Очистка ресурсов context.WithTimeout
	defer cancel()

	// Остановка сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	fmt.Println("Сервер успешно остановлен")
}
