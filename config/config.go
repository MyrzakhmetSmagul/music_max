package config

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func loadEnvFile(filename string) error {
	// Открываем файл .env
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Сканируем файл построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Пропускаем пустые строки и комментарии
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Разделяем строку на ключ и значение по символу "="
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			slog.Warn("loadENVFile invalid line: ", slog.Any("line", line))
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Устанавливаем переменную окружения
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func InitConfig(filename string) error {
	slog.Debug("init configurations")
	err := loadEnvFile(filename)
	if err != nil {
		err = fmt.Errorf("error while loading env file: %w", err)
		return err
	}

	return nil
}
