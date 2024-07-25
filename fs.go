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

// FileInfoWithSize - это структура, которая комбинирует информацию о файле (os.FileInfo)
// с дополнительным полем для хранения размера файла или директории.
type FileInfoWithSize struct {
	// Информация об имени файла
	os.FileInfo
	// Информация о размере файл
	Size int64
}

func main() {
	start := time.Now()

	root, sortOrder := parseFlags()

	// Получения списка файлов из каталога
	files, err := ioutil.ReadDir(*root)
	if err != nil {
		fmt.Printf("Ошибка чтения директория: %v\n", err)
		return
	}
	// Вывод списка файлов из каталога
	processFiles(*root, files)

	// Обработка файлов и директорий для получения размеров
	fileInfoWithSizes := processFiles(*root, files)

	// Сортировка файлов и директорий по размеру
	sortFiles(fileInfoWithSizes, *sortOrder)

	// Вывод результатов в виде таблицы
	fmt.Printf("%-10s %-30s %15s\n", "Тип", "Имя", "Размер")
	for _, fileInfo := range fileInfoWithSizes {
		// Определение типа: файл или директория
		fileType := "Файл"
		if fileInfo.IsDir() {
			fileType = "Дир"
		}
		// Вывод информации о файле или директории
		fmt.Printf("%-10s %-20s %15s\n", fileType, fileInfo.Name(), formatSize(fileInfo.Size))
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
		// file.Join - предназначенна для безопасного объединения частей пути в один путь. Она учитывает особенности операционной системы,
		// такие как правильное использование разделителей путей
		fullPath := filepath.Join(root, file.Name())
		size := file.Size()
		if file.IsDir() {
			// Если это директория, вычисляем её размер
			size = getDirSize(fullPath)
		}
		// Добавление информации о файле или директории в результат
		result = append(result, FileInfoWithSize{file, size})
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
		fmt.Printf("Ошибка чтения директория в рекурсивной функции: %v\n", err)
	}
	return size
}

// Функция для сортировки файлов и директорий по размеру
func sortFiles(files []FileInfoWithSize, order string) {
	sort.Slice(files, func(i, j int) bool {
		if order == "asc" {
			// Сортировка по возрастанию
			// Если размер первого файла (files[i].Size) меньше размера второго файла (files[j].Size), функция возвращает true
			// Это означает, что первый файл должен быть перед вторым в отсортированном списке
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
	// size - размер файла или директории в байтах
	// div - переменная для хранения текущего масштаба
	// exp - переменная для хранения экспоненты, определяющей единицу измерения
	// (0 для байтов, 1 для килобайтов, 2 для мегабайтов и т.д.)
	div, exp := int64(unit), 0
	// Цикл с делением размера файла на единицу измерения (1024)
	for n := size / unit; n >= unit; n /= unit {
		// Увеличиваем масштаб, умножая div на единицу измерения (1024)
		// Это переводит масштаб на следующий уровень (байты -> килобайты -> мегабайты и т.д.)
		div *= unit

		// Увеличиваем экспоненту, чтобы указать на следующую единицу измерения
		// exp = 0 для байтов, 1 для килобайтов, 2 для мегабайтов и т.д.
		exp++
	}
	// Форматированный вывод с одной цифрой после запятой
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
