package models

import "time"

// Habit представляет привычку, которую отслеживает пользователь.
type Habit struct {
    Name         string    // Название привычки
    Description  string    // Описание привычки (опционально)
    StartDate    time.Time // Дата начала привычки
    Completed    []time.Time // Массив дат, когда привычка была выполнена
}
