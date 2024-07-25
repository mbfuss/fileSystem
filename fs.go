package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// FileInfoWithSize - структура, которая комбинирует информацию о файле с дополнительным полем для хранения размера файла или директории.
type FileInfoWithSize struct {
	Name   string // Имя файла
	IsFile bool   // Тип файла
	Size   int64  // Размер файла
}

func main() {
	start := time.Now()
	// Создание флагов и получение их значений
	root, sortOrder := parseFlags()

	// Получение списка файлов из каталога
	entries, err := os.ReadDir(*root)
	if err != nil {
		fmt.Printf("Ошибка чтения директории: %v\n", err)
		return
	}

	// Обработка файлов и директорий для получения размеров
	fileInfoWithSizes := processFiles(*root, entries)

	// Сортировка файлов и директорий по размеру
	sortFiles(fileInfoWithSizes, *sortOrder)

	// Вывод результатов в виде таблицы
	printFileInfoTable(fileInfoWithSizes)

	duration := time.Since(start)
	fmt.Println("Время выполнения программы", duration)
}

// parseFlags - функция для создания флагов и проверки их на валидность
func parseFlags() (root *string, sortOrder *string) {
	root = flag.String("root", "", "Путь до корневой директории")
	sortOrder = flag.String("sort", "asc", "Порядок сортировки: asc (возрастание) или des (убывание)")
	flag.Parse()

	if *root == "" || (*sortOrder != "asc" && *sortOrder != "des") {
		flag.Usage()
		log.Fatal("Флаги переданы неверно")
	}

	return root, sortOrder
}

// processFiles - принимает корневую директорию и список файлов/директорий, вычисляет размер каждого элемента,
// и возвращает список структур FileInfoWithSize, которые содержат информацию о файлах/директориях и их размерах.
func processFiles(root string, entries []os.DirEntry) []FileInfoWithSize {
	var result []FileInfoWithSize
	// Создается полный путь к текущему элементу, используя функцию filepath.Join, которая корректно объединяет корневую директорию (root)
	// и имя текущего элемента (entry.Name()).
	for _, entry := range entries {
		// Полный путь к файлу или директории
		fullPath := filepath.Join(root, entry.Name())
		// Получаем информацию о каждом файле
		fileInfo, err := entry.Info()
		if err != nil {
			fmt.Printf("Ошибка получения информации о файле: %v\n", err)
			continue
		}
		// Получаем размер текущего элемента
		size := fileInfo.Size()
		// Если это не директория, то это файл
		isFile := !entry.IsDir()
		if !isFile {
			// Если это директория, вычисляем её размер
			size, err = getDirSize(fullPath)
			if err != nil {
				fmt.Printf("Ошибка чтения директории в рекурсивной функции: %v\n", err)
			}
		}
		// Добавление информации о файле или директории в результат
		result = append(result, FileInfoWithSize{Name: entry.Name(), IsFile: isFile, Size: size})
	}
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

// Функция для вывода информации о файлах и директориях в виде таблицы
func printFileInfoTable(fileInfoWithSizes []FileInfoWithSize) {
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
}
