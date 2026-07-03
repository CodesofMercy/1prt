package models

import "time"

// Habit представляет привычку, которую отслеживает пользователь.
type Habit struct {
	Name        string      `json:"name"`                  // Название привычки
	Description string      `json:"description,omitempty"` // Описание привычки (опционально)
	StartDate   time.Time   `json:"start_date"`            // Дата начала привычки
	Completed   []time.Time `json:"completed"`             // Даты, когда привычка была выполнена
}

// sameDay сообщает, приходятся ли два момента времени на один календарный день.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// IsDoneOn сообщает, была ли привычка выполнена в день, на который приходится t.
func (h *Habit) IsDoneOn(t time.Time) bool {
	for _, c := range h.Completed {
		if sameDay(c, t) {
			return true
		}
	}
	return false
}

// DoneToday сообщает, отмечена ли привычка выполненной сегодня.
func (h *Habit) DoneToday() bool {
	return h.IsDoneOn(time.Now())
}

// MarkDone отмечает привычку выполненной в день t.
// Повторная отметка в тот же день игнорируется; возвращает true, если отметка добавлена.
func (h *Habit) MarkDone(t time.Time) bool {
	if h.IsDoneOn(t) {
		return false
	}
	h.Completed = append(h.Completed, t)
	return true
}

// Streak возвращает текущую серию: сколько дней подряд привычка выполнялась.
// Серия не считается прерванной, пока текущий день ещё не закончился:
// если сегодня отметки нет, отсчёт начинается со вчерашнего дня.
func (h *Habit) Streak(now time.Time) int {
	day := now
	if !h.IsDoneOn(day) {
		day = day.AddDate(0, 0, -1)
	}
	streak := 0
	for h.IsDoneOn(day) {
		streak++
		day = day.AddDate(0, 0, -1)
	}
	return streak
}
