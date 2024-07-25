package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// FileInfoWithSize - это структура, которая комбинирует информацию о файле с дополнительным полем для хранения размера файла или директории.
type FileInfoWithSize struct {
	Name   string
	IsFile bool
	Size   int64
}

func main() {
	start := time.Now()

	root, sortOrder := parseFlags()

	// Получение списка файлов из каталога
	files, err := ioutil.ReadDir(*root)
	if err != nil {
		fmt.Printf("Ошибка чтения директории: %v\n", err)
		return
	}

	// Обработка файлов и директорий для получения размеров
	fileInfoWithSizes := processFiles(*root, files)

	// Сортировка файлов и директорий по размеру
	sortFiles(fileInfoWithSizes, *sortOrder)

	// Вывод результатов в виде таблицы
	fmt.Printf("%-10s %-30s %15s\n", "Тип", "Имя", "Размер")
	for _, fileInfo := range fileInfoWithSizes {
		// Определение типа: файл или директория
		fileType := "Файл"
		if !fileInfo.IsFile {
			fileType = "Дир"
		}
		// Вывод информации о файле или директории
		fmt.Printf("%-10s %-30s %15s\n", fileType, fileInfo.Name, formatSize(fileInfo.Size))
	}

	duration := time.Since(start)
	fmt.Println("Время выполнения программы", duration)
}

// parseFlags - функция для создания флагов и проверки их на валидность
func parseFlags() (root *string, sortOrder *string) {
	root = flag.String("root", "", "Путь до корневой директории")
	sortOrder = flag.String("sort", "", "Порядок сортировки: asc (возрастание) или des (убывание)")
	flag.Parse()

	if *root == "" || *sortOrder == "" {
		flag.Usage()
		log.Fatal("Флаги переданы неверно")
	}

	return root, sortOrder
}

// processFiles - принимает корневую директорию и список файлов/директорий, вычисляет размер каждого элемента,
// и возвращает список структур FileInfoWithSize, которые содержат информацию о файлах/директориях и их размерах.
func processFiles(root string, files []os.FileInfo) []FileInfoWithSize {
	var result []FileInfoWithSize
	for _, file := range files {
		// Полный путь к файлу или директории
		fullPath := filepath.Join(root, file.Name())
		size := file.Size()
		isFile := !file.IsDir()
		if !isFile {
			// Если это директория, вычисляем её размер
			size = getDirSize(fullPath)
		}
		// Добавление информации о файле или директории в результат
		result = append(result, FileInfoWithSize{Name: file.Name(), IsFile: isFile, Size: size})
	}
	return result
}

// getDirSize - функция которая вычисляет размер директории
func getDirSize(path string) int64 {
	var size int64
	// Рекурсивно проходит по всем файлам и поддиректориям, начиная с указанного пути (path)
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Если это файл, добавляем его размер к общему размеру
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Ошибка чтения директории в рекурсивной функции: %v\n", err)
	}
	return size
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
	const unit = 1024
	if size < unit {
		// Если размер меньше 1024 байт, выводим в байтах
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	// Цикл с делением размера файла на единицу измерения (1024)
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	// Форматированный вывод с одной цифрой после запятой
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
