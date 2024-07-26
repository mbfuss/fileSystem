package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// FileInfoWithSize - структура, которая комбинирует информацию о файле с дополнительным полем для хранения размера файла или директории.
type FileInfoWithSize struct {
	Name   string `json:"name"`    // Имя файла
	IsFile bool   `json:"is_file"` // Тип файла
	Size   int64  `json:"size"`    // Размер файла
}

const desc = "desc"
const asc = "asc"

func main() {
	// Регистрация обработчика для пути /fs:
	http.HandleFunc("/fs", handleFileRequest)
	fmt.Println("Сервер запущен на порту 80")
	// Запускает сервера, который будет прослушивать порт 80
	// nil -- использовать глобальный маршрутизатор http.DefaultServeMux для обработки запросов.
	log.Fatal(http.ListenAndServe(":80", nil))
}

// handleFileRequest - функция, которая обрабатывает http запросы
// http.ResponseWriter — интерфейс, который предоставляет методы для формирования и отправки HTTP-ответов клиенту
// http.Request — это указатель на структуру http.Request в Go, которая представляет собой HTTP-запрос, отправленный клиентом на сервер
func handleFileRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Получение параметров запроса
	root := r.URL.Query().Get("root")
	sortOrder := r.URL.Query().Get("sort")

	// Проверка параметров запроса
	if root == "" || (sortOrder != asc && sortOrder != desc) {
		http.Error(w, "Неверные параметры запроса", http.StatusBadRequest)
		return
	}

	// Получение списка файлов из каталога
	entries, err := os.ReadDir(root)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка чтения директории: %v", err), http.StatusInternalServerError)
		return
	}

	// Обработка файлов и директорий для получения размеров
	fileInfoWithSizes := processFiles(root, entries)

	// Сортировка файлов и директорий по размеру
	sortFiles(fileInfoWithSizes, sortOrder)

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

// processFiles - принимает корневую директорию и список файлов/директорий, вычисляет размер каждого элемента,
// и возвращает список структур FileInfoWithSize, которые содержат информацию о файлах/директориях и их размерах.
func processFiles(root string, entries []os.DirEntry) []FileInfoWithSize {
	var wg sync.WaitGroup
	var result = make([]FileInfoWithSize, len(entries))
	for i, entry := range entries {
		wg.Add(1)
		go func(i int, entry os.DirEntry, wg *sync.WaitGroup, result []FileInfoWithSize) {
			defer wg.Done()
			fullPath := filepath.Join(root, entry.Name())
			fileInfo, err := entry.Info()
			if err != nil {
				fmt.Printf("Ошибка получения информации о файле: %v\n", err)
				return
			}
			size := fileInfo.Size()
			isFile := !entry.IsDir()
			if !isFile {
				size, err = getDirSize(fullPath)
				if err != nil {
					fmt.Printf("Ошибка чтения директории в рекурсивной функции: %v\n", err)
				}
			}
			result[i] = FileInfoWithSize{Name: entry.Name(), IsFile: isFile, Size: size}
		}(i, entry, &wg, result)
	}
	wg.Wait()
	return result
}

// getDirSize - функция, которая вычисляет размер директории
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		if info.IsDir() {
			size += 4000
		}
		return nil
	})

	return size, err
}

// sortFiles - функция для сортировки файлов и директорий по размеру
func sortFiles(files []FileInfoWithSize, order string) {
	sort.Slice(files, func(i, j int) bool {
		if order == "asc" {
			return files[i].Size < files[j].Size
		}
		return files[i].Size > files[j].Size
	})
}
