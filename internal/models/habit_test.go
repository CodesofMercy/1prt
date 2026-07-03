package models

import (
	"testing"
	"time"
)

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 12, 0, 0, 0, time.Local)
}

func TestMarkDoneIdempotentPerDay(t *testing.T) {
	h := &Habit{Name: "Чтение"}
	now := day(2026, 7, 3)

	if !h.MarkDone(now) {
		t.Fatal("первая отметка за день должна добавляться")
	}
	if h.MarkDone(now.Add(2 * time.Hour)) {
		t.Fatal("повторная отметка в тот же день не должна добавляться")
	}
	if len(h.Completed) != 1 {
		t.Fatalf("ожидалась 1 отметка, получено %d", len(h.Completed))
	}
	if !h.MarkDone(now.AddDate(0, 0, 1)) {
		t.Fatal("отметка на следующий день должна добавляться")
	}
}

func TestIsDoneOn(t *testing.T) {
	h := &Habit{Completed: []time.Time{day(2026, 7, 1)}}
	if !h.IsDoneOn(day(2026, 7, 1).Add(5 * time.Hour)) {
		t.Error("должна быть выполнена в тот же день независимо от времени суток")
	}
	if h.IsDoneOn(day(2026, 7, 2)) {
		t.Error("не должна быть выполнена в другой день")
	}
}

func TestStreak(t *testing.T) {
	now := day(2026, 7, 3)

	tests := []struct {
		name string
		done []time.Time
		want int
	}{
		{"нет отметок", nil, 0},
		{"только сегодня", []time.Time{now}, 1},
		{"три дня подряд включая сегодня", []time.Time{now, now.AddDate(0, 0, -1), now.AddDate(0, 0, -2)}, 3},
		{"сегодня ещё нет, вчера и позавчера есть", []time.Time{now.AddDate(0, 0, -1), now.AddDate(0, 0, -2)}, 2},
		{"серия прервана позавчера", []time.Time{now.AddDate(0, 0, -2), now.AddDate(0, 0, -3)}, 0},
		{"разрыв внутри серии", []time.Time{now, now.AddDate(0, 0, -2)}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Habit{Completed: tt.done}
			if got := h.Streak(now); got != tt.want {
				t.Errorf("Streak() = %d, ожидалось %d", got, tt.want)
			}
		})
	}
}
