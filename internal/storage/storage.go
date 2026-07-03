// Package storage отвечает за сохранение и загрузку привычек на диск в формате JSON.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/CodesofMercy/1prt/internal/models"
)

// DefaultPath возвращает путь к файлу данных в каталоге конфигурации пользователя,
// например ~/.config/1prt/habits.json (Linux) или %AppData%\1prt\habits.json (Windows).
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("не удалось определить каталог конфигурации: %w", err)
	}
	return filepath.Join(dir, "1prt", "habits.json"), nil
}

// Load читает привычки из файла. Если файла ещё нет, возвращает пустой список.
func Load(path string) ([]models.Habit, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return []models.Habit{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать %s: %w", path, err)
	}
	var habits []models.Habit
	if err := json.Unmarshal(data, &habits); err != nil {
		return nil, fmt.Errorf("повреждён файл данных %s: %w", path, err)
	}
	return habits, nil
}

// Save записывает привычки в файл атомарно: сначала во временный файл,
// затем переименовывает, чтобы при сбое не потерять предыдущие данные.
func Save(path string, habits []models.Habit) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("не удалось создать каталог данных: %w", err)
	}
	data, err := json.MarshalIndent(habits, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать данные: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("не удалось записать %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("не удалось сохранить %s: %w", path, err)
	}
	return nil
}
