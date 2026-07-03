package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CodesofMercy/1prt/internal/models"
)

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "habits.json")
	habits, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(habits) != 0 {
		t.Fatalf("ожидался пустой список, получено %d", len(habits))
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "habits.json")
	want := []models.Habit{
		{
			Name:      "Зарядка",
			StartDate: time.Date(2026, 7, 1, 8, 0, 0, 0, time.UTC),
			Completed: []time.Time{time.Date(2026, 7, 2, 9, 30, 0, 0, time.UTC)},
		},
		{Name: "Чтение", StartDate: time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC)},
	}

	if err := Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("ожидалось %d привычек, получено %d", len(want), len(got))
	}
	for i := range want {
		if got[i].Name != want[i].Name {
			t.Errorf("привычка %d: имя %q, ожидалось %q", i, got[i].Name, want[i].Name)
		}
		if !got[i].StartDate.Equal(want[i].StartDate) {
			t.Errorf("привычка %d: дата начала не совпала", i)
		}
		if len(got[i].Completed) != len(want[i].Completed) {
			t.Errorf("привычка %d: отметок %d, ожидалось %d", i, len(got[i].Completed), len(want[i].Completed))
		}
	}
}

func TestLoadCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "habits.json")
	if err := Save(path, nil); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// Портим файл вручную.
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("ожидалась ошибка на повреждённом файле")
	}
}
