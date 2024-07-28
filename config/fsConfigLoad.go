package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv - чтение файла .env
func LoadEnv(filename string) error {
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
