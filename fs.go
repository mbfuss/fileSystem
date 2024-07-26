package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// FileInfoWithSize - структура, которая комбинирует информацию о файле с дополнительным полем для хранения размера файла или директории.
type FileInfoWithSize struct {
	Name   string `json:"name"`    // Имя файла
	IsFile bool   `json:"is_file"` // Тип файла
	Size   int64  `json:"size"`    // Размер файла
}

type MapperFileInfo struct {
	Name   string `json:"name"`    // Имя файла
	IsFile string `json:"is_file"` // Тип файла
	Size   string `json:"size"`    // Размер файла
}

const desc = "desc"
const asc = "asc"

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue // Пропустить пустые строки и комментарии
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Пропустить строки без "="
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value) // Установить переменную окружения
	}

	return scanner.Err()
}

func main() {
	// Загружаем переменные из .env файла
	err := loadEnv("serverPort.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Читаем порт из переменной окружения
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "80" // Значение по умолчанию, если PORT не задан
	}
	// Регистрация обработчика для пути /fs
	http.HandleFunc("/fs", handleFileRequest)
	fmt.Printf("Сервер запущен на порту %s\n", port)
	// Запускает сервера, который будет прослушивать порт 80
	// nil -- использовать глобальный маршрутизатор http.DefaultServeMux для обработки запросов.
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
	err = json.NewEncoder(w).Encode(mapFileInfoWithSizeToMapperFileInfo(fileInfoWithSizes))
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
			// Создается полный путь к текущему элементу, используя функцию filepath.Join, которая корректно объединяет корневую директорию (root)
			// и имя текущего элемента (entry.Name()).
			fullPath := filepath.Join(root, entry.Name())
			// Получаем информацию о каждом файле
			fileInfo, err := entry.Info()
			if err != nil {
				fmt.Printf("Ошибка получения информации о файле: %v\n", err)
				return
			}
			// Получаем размер текущего элемента
			size := fileInfo.Size()
			isFile := !entry.IsDir()
			if !isFile {
				size, err = getDirSize(fullPath)
				if err != nil {
					fmt.Printf("Ошибка чтения директории в рекурсивной функции: %v\n", err)
				}
			}
			result[i] = FileInfoWithSize{
				Name:   entry.Name(),
				IsFile: isFile,
				Size:   size,
			}
		}(i, entry, &wg, result)
	}
	wg.Wait()
	return result
}

// getDirSize - функция которая вычисляет размер директории
func getDirSize(path string) (int64, error) {
	var size int64
	// Рекурсивно проходит по всем файлам и поддиректориям, начиная с указанного пути (path)
	// Анонимная функция вызывается для каждого файла и директории, найденных во время обхода
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Если это файл, добавляем его размер к общему размеру
		if !info.IsDir() {
			size += info.Size()
		}
		// Учет размера папки
		if info.IsDir() {
			size += 4000
		}

		return nil
	})

	return size, err
}

// Функция для сортировки файлов и директорий по размеру
func sortFiles(files []FileInfoWithSize, order string) {
	sort.Slice(files, func(i, j int) bool {
		if order == "asc" {
			// Сортировка по возрастанию
			return files[i].Size < files[j].Size
		}
		// Сортировка по убыванию
		return files[i].Size > files[j].Size
	})
}

// Функция для форматирования размера файла или директории в читаемый вид
func formatSize(size int64) string {
	const unit = 1000
	if size < unit {
		// Если размер меньше 1 килобайта, выводим в байтах
		return fmt.Sprintf("%d B", size)
	}
	// div: переменная для хранения текущего масштаба единицы измерения. Изначально устанавливается в 1000 (1 килобайт)
	// exp: переменная для хранения экспоненты, указывающей на текущую единицу измерения. Изначально установлена в 0 (для байт)
	div, exp := int64(unit), 0
	// Цикл с делением размера файла на единицу измерения (1024)
	for n := size / unit; n >= unit; {
		n /= unit
		div *= unit
		exp++
	}
	// Форматированный вывод с одной цифрой после запятой
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGT"[exp])
}

func formatIsFile(fileType bool) string {
	if fileType {
		return "Файл"
	} else {
		return "Дир"
	}

}

// mapFileInfoWithSizeToMapperFileInfo - функция для преобразования среза FileInfoWithSize в срез MapperFileInfo для корректного отображения в json
func mapFileInfoWithSizeToMapperFileInfo(fileInfos []FileInfoWithSize) []MapperFileInfo {
	var mappedInfos []MapperFileInfo
	for _, fileInfo := range fileInfos {

		mappedInfo := MapperFileInfo{
			Name:   fileInfo.Name,
			IsFile: formatIsFile(fileInfo.IsFile),
			Size:   formatSize(fileInfo.Size), // Применяем форматирование размера
		}
		mappedInfos = append(mappedInfos, mappedInfo)
	}
	return mappedInfos
}
