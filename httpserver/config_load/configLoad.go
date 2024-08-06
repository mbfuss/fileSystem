package config_load

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv - чтение файла .env
func LoadEnv(filename string) error {
	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем сканер для чтения file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Считываем построчно
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue // Пропустить пустые строки и комментарии
		}
		// key - SERVER_PORT, value - SERVER_PORT
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Пропустить строки без "="
		}

		// .TrimSpace - удаляет табуляции и символы новой строки
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value) // Установить переменную окружения
	}

	return scanner.Err()
}
