package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// FileInfoWithSize - структура, которая комбинирует информацию о файле с дополнительным полем для хранения размера файла или директории
type FileInfoWithSize struct {
	Name       string `json:"name"`        // Имя файла
	IsFile     string `json:"is_file"`     // Тип файла
	Size       int64  `json:"-"`           // Размер файла
	FormatSize string `json:"format_size"` // Размер файла после фармотирования
}

// ResponseData - структура для хранения ответа сервера
type ResponseData struct {
	ErrorCode    int                `json:"error_code"`
	ErrorMessage string             `json:"error_message"`
	Data         []FileInfoWithSize `json:"data"`
	Root         string             `json:"root"`
}

const (
	desc = "desc"
	asc  = "asc"
)

// GetRequestParams - получает и проверяет параметры запроса root и sortOrder.
// На вход: r *http.Request: Указатель на объект http.Request, который представляет собой HTTP-запрос. Этот объект содержит все данные о запросе, включая URL, заголовки, тело запроса и параметры запроса.
func GetRequestParams(r *http.Request) (string, string, error) {
	root := r.URL.Query().Get("root")
	sortOrder := r.URL.Query().Get("sort")

	if root == "" || (sortOrder != asc && sortOrder != desc) {
		return "", "", errors.New("неверные параметры запроса")
	}

	return root, sortOrder, nil
}

// ProcessFiles - принимает корневую директорию и список файлов/директорий, вычисляет размер каждого элемента,
// и возвращает список структур FileInfoWithSize, которые содержат информацию о файлах/директориях и их размерах
func ProcessFiles(root string, entries []os.DirEntry) []FileInfoWithSize {
	var wg sync.WaitGroup
	var result = make([]FileInfoWithSize, len(entries))
	for i, entry := range entries {
		wg.Add(1)
		go func(i int, entry os.DirEntry, wg *sync.WaitGroup, result []FileInfoWithSize) {
			defer wg.Done()
			// Создается полный путь к текущему элементу, используя функцию filepath.Join, которая корректно объединяет корневую директорию (root)
			// и имя текущего элемента (entry.Name())
			fullPath := filepath.Join(root, entry.Name())
			// Получаем информацию о каждом файле
			fileInfo, err := entry.Info()
			if err != nil {
				fmt.Printf("Ошибка получения информации о файле: %v\n", err)
				return
			}
			// Получаем размер текущего элемента
			size := fileInfo.Size()
			// Если не директория, то файл
			isFile := !entry.IsDir()
			if !isFile {
				size, err = getDirSize(fullPath)
				if err != nil {
					fmt.Printf("Ошибка чтения директории в рекурсивной функции: %v\n", err)
				}
			}
			result[i] = FileInfoWithSize{
				Name:       entry.Name(),
				IsFile:     formatIsFile(isFile),
				Size:       size,
				FormatSize: formatSize(size),
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
		size += info.Size()
		// Учет размера папки
		if info.IsDir() {
			size += 4000
		}

		return nil
	})

	return size, err
}

// SortFiles - функция для сортировки файлов и директорий по размеру
func SortFiles(files []FileInfoWithSize, order string) {
	// i и j, указывающие на элементы в срезе, и возвращает булево значение,
	// указывающее, должен ли элемент с индексом i быть перед элементом с индексом j
	sort.Slice(files, func(i, j int) bool {
		if order == asc {
			// Сортировка по возрастанию
			return files[i].Size < files[j].Size
		}
		// Сортировка по убыванию
		return files[i].Size > files[j].Size
	})
}

// formatSize - функция для форматирования размера файла или директории в читаемый вид
func formatSize(size int64) string {
	const unit = 1000
	if size < unit {
		// Если размер меньше 1 килобайта, выводим в байтах
		return fmt.Sprintf("%d B", size)
	}
	// div: переменная для хранения текущего масштаба единицы измерения. Изначально устанавливается в 1000 (1 килобайт)
	// exp: переменная для хранения экспоненты, указывающей на текущую единицу измерения. Изначально установлена в 0 (для байт)
	div, exp := int64(unit), 0
	// Цикл с делением размера файла на единицу измерения (1000)
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
