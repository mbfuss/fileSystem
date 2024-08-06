package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mbfuss/sortingFiles/httpserver/configLoad"
	"github.com/mbfuss/sortingFiles/httpserver/service"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Время для завершения сервера
const (
	timeOut    = 5 * time.Second
	serverPort = "SERVER_PORT"
	rootDir    = "ROOT_DIR"
)

// HandleFileRequest - функция, которая обрабатывает http запросы
// http.ResponseWriter — интерфейс, который предоставляет методы для формирования и отправки HTTP-ответов клиенту
// http.Request — это указатель на структуру http.Request в Go, которая представляет собой HTTP-запрос, отправленный клиентом на сервер
func HandleFileRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Получение параметров запроса
	root, sortOrder, err := service.GetRequestParams(r)
	rootDir := os.Getenv(rootDir)
	if err != nil {
		response := service.ResponseData{
			ErrorCode:    1,
			ErrorMessage: fmt.Sprintf("Ошибка получения пути, будет выбран путь по умолчанию:", err),
			Data:         nil,
			Root:         rootDir,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Получение списка файлов из каталога
	entries, err := os.ReadDir(root)
	if err != nil {
		response := service.ResponseData{
			ErrorCode:    1,
			ErrorMessage: fmt.Sprintf("Ошибка чтения директории: %v", err),
			Data:         nil,
			Root:         root,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Обработка файлов и директорий для получения размеров
	fileInfoWithSizes := service.ProcessFiles(root, entries)
	// Сортировка файлов и директорий по размеру
	service.SortFiles(fileInfoWithSizes, sortOrder)
	// Тип отправки ответа на клиент
	response := service.ResponseData{
		ErrorCode:    0,
		ErrorMessage: "",
		Data:         fileInfoWithSizes,
		Root:         rootDir,
	}

	// Используется для формирования и отправки HTTP-ответов клиенту в формате json
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w) создает новый JSON-энкодер, который будет записывать данные прямо в http.ResponseWriter (w), то есть отправлять данные клиенту
	// .Encode(fileInfoWithSizes) преобразует fileInfoWithSizes (срез структур FileInfoWithSize) в JSON-формат и отправляет его в ответе
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка кодирования в JSON: %v", err), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	fmt.Printf("Обработка запроса: %v\n", duration)

	// ЛОГИКА ОТПРАВКИ POST ЗАПРОСА НА APACHE PHP

	// Инициализируем переменную для хранения общего размера файлов
	totalSize := int64(0)
	// Проходим по всем файлам и суммируем их размеры
	for _, fileInfo := range fileInfoWithSizes {
		totalSize += fileInfo.Size
	}

	// Подготавливаем данные для отправки в формате JSON
	logData := map[string]interface{}{
		"path":     root,                    // Путь к директории
		"size":     totalSize,               // Общий размер всех файлов
		"duration": duration.Milliseconds(), // Длительность обработки запроса в миллисекундах
	}

	// Кодируем данные в формат JSON
	logDataJson, err := json.Marshal(logData)
	if err != nil {
		// Если произошла ошибка при кодировании данных в JSON, выводим сообщение об ошибке и завершаем функцию
		fmt.Printf("Ошибка кодирования данных для логирования: %v\n", err)
		return
	}

	// Определяем URL для отправки POST-запроса на PHP-скрипт
	aRoot := os.Getenv("APACHE_ROOT")
	phpUrl := aRoot + "/setStat.php"
	// Отправляем POST-запрос с данными в формате JSON
	resp, err := http.Post(phpUrl, "application/json", bytes.NewBuffer(logDataJson))
	if err != nil {
		// Если произошла ошибка при отправке запроса, выводим сообщение об ошибке и завершаем функцию
		fmt.Printf("Ошибка отправки данных на сервер Apache PHP: %v\n", err)
		return
	}

	// Закрываем тело ответа после завершения работы с ним
	defer resp.Body.Close()

	// Проверяем статус ответа от PHP-сервера
	if resp.StatusCode != http.StatusOK {
		// Если статус ответа не 200 OK, выводим сообщение об ошибке
		fmt.Printf("Сервер Apache PHP вернул ошибку: %v\n", resp.Status)
	}

}

// HandleGetFileInfo - обработчик для получения данных из MySQL через PHP
func HandleGetFileInfo(w http.ResponseWriter, _ *http.Request) {
	aRoot := os.Getenv("APACHE_ROOT")
	phpUrl := aRoot + "/getStat.php"
	// Отправка GET-запроса к PHP-скрипту
	resp, err := http.Get(phpUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при отправке запроса на PHP сервер: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверка на успешный ответ
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("PHP сервер вернул ошибку: %v", resp.Status), http.StatusInternalServerError)
		return
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при чтении ответа от PHP сервера: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправка HTML-данных клиенту
	w.Header().Set("Content-Type", "text/html")
	w.Write(body)
}

func StartApp() {

	// Загружаем переменные из .env файла
	err := configLoad.LoadEnv("config/serverPort.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Читаем порт из переменной окружения
	port := os.Getenv(serverPort)

	// Создание нового сервера
	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	// Регистрация обработчиков
	http.HandleFunc("/fs", HandleFileRequest)
	// Отображение представления localhost:SERVER_PORT
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)
	http.HandleFunc("/getfileinfo", HandleGetFileInfo)

	// Канал для получения системных сигналов
	stop := make(chan os.Signal, 1)
	// Функция signal.Notify регистрирует канал stop для получения уведомлений о сигнале os.Interrupt (Ctrl+C) и os.Kill (сигнал принудительного завершения процесса)
	// При получении одного из этих сигналов, он будет отправлен в канал stop
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Запуск сервера в отдельной горутине
	go func() {
		fmt.Printf("Сервер запущен на порту %s\n", port)
		// Запуск сервера
		err := server.ListenAndServe()
		// Если ошибка не пустая, и если сервер не был корректно остановлен
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Сервер остановлен с ошибкой: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	fmt.Println("Получен сигнал завершения, остановка сервера...")

	//  Если сервер не завершится в течение 5 секунд, сервер принудительно завершит работу
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	// Используется для отложенного вызова функции cancel. Это значит, что функция cancel будет вызвана автоматически,
	// когда выполнение функции ServerStart завершится, даже если произойдёт ошибка. Очистка ресурсов context.WithTimeout
	defer cancel()

	// Остановка сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	fmt.Println("Сервер успешно остановлен")
}
