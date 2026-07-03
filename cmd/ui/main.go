package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/CodesofMercy/1prt/internal/models"
	"github.com/CodesofMercy/1prt/internal/storage"
)

const appVersion = "0.1.0"

// tracker хранит состояние приложения: список привычек и путь к файлу данных.
type tracker struct {
	habits   []models.Habit
	dataPath string
	window   fyne.Window
	status   *widget.Label
	list     *fyne.Container
}

// capitalize делает первую букву заглавной с учётом кириллицы и других алфавитов.
func capitalize(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func (t *tracker) habitExists(name string) bool {
	for _, h := range t.habits {
		if strings.EqualFold(h.Name, name) {
			return true
		}
	}
	return false
}

// save записывает привычки на диск; ошибку показывает пользователю.
func (t *tracker) save() {
	if err := storage.Save(t.dataPath, t.habits); err != nil {
		dialog.ShowError(err, t.window)
	}
}

func (t *tracker) addHabit(name string) {
	name = capitalize(strings.TrimSpace(name))
	if name == "" {
		return
	}
	if t.habitExists(name) {
		dialog.ShowInformation("Уже есть", "Привычка «"+name+"» уже добавлена.", t.window)
		return
	}
	t.habits = append(t.habits, models.Habit{
		Name:      name,
		StartDate: time.Now(),
		Completed: []time.Time{},
	})
	t.save()
	t.refresh()
	t.status.SetText("Привычка добавлена: " + name)
}

func (t *tracker) markDone(index int) {
	if index < 0 || index >= len(t.habits) {
		return
	}
	h := &t.habits[index]
	if h.MarkDone(time.Now()) {
		t.save()
		t.refresh()
		t.status.SetText("Отмечено на сегодня: " + h.Name)
	}
}

func (t *tracker) deleteHabit(index int) {
	if index < 0 || index >= len(t.habits) {
		return
	}
	name := t.habits[index].Name
	dialog.ShowConfirm("Удалить привычку",
		"Удалить «"+name+"» вместе со всей историей отметок?",
		func(ok bool) {
			if !ok {
				return
			}
			t.habits = append(t.habits[:index], t.habits[index+1:]...)
			t.save()
			t.refresh()
			t.status.SetText("Привычка удалена: " + name)
		}, t.window)
}

// refresh перестраивает список привычек на экране.
func (t *tracker) refresh() {
	t.list.Objects = nil
	if len(t.habits) == 0 {
		empty := widget.NewLabel("Пока нет ни одной привычки — добавьте первую!")
		t.list.Add(empty)
		t.list.Refresh()
		return
	}

	now := time.Now()
	for i, h := range t.habits {
		index := i

		nameLabel := widget.NewLabel(h.Name)
		nameLabel.TextStyle = fyne.TextStyle{Bold: true}

		streak := h.Streak(now)
		streakLabel := widget.NewLabel(fmt.Sprintf("серия: %d дн.", streak))

		var doneButton *widget.Button
		if h.IsDoneOn(now) {
			doneButton = widget.NewButtonWithIcon("Выполнено сегодня", theme.ConfirmIcon(), nil)
			doneButton.Disable()
		} else {
			doneButton = widget.NewButton("Отметить на сегодня", func() {
				t.markDone(index)
			})
		}

		deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			t.deleteHabit(index)
		})

		row := container.NewBorder(nil, nil,
			container.NewHBox(nameLabel, streakLabel),   // слева
			container.NewHBox(doneButton, deleteButton), // справа
		)
		t.list.Add(row)
		t.list.Add(widget.NewSeparator())
	}
	t.list.Refresh()
}

// showAddDialog открывает диалог добавления привычки.
// Добавление выполняется один раз — либо по Enter, либо по кнопке «Добавить».
func (t *tracker) showAddDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Например: читать 10 страниц")

	form := dialog.NewForm("Новая привычка", "Добавить", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Название", nameEntry)},
		func(ok bool) {
			if ok {
				t.addHabit(nameEntry.Text)
			}
		}, t.window)

	nameEntry.OnSubmitted = func(text string) {
		form.Hide() // закрываем диалог, добавляем один раз
		t.addHabit(text)
	}

	form.Resize(fyne.NewSize(400, form.MinSize().Height))
	form.Show()
	t.window.Canvas().Focus(nameEntry)
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("1% Habit Tracker v" + appVersion)
	window.Resize(fyne.NewSize(800, 600))

	dataPath, err := storage.DefaultPath()
	if err != nil {
		log.Fatal(err)
	}

	t := &tracker{
		dataPath: dataPath,
		window:   window,
		status:   widget.NewLabel("Добро пожаловать в приложение 1%! Каждый день делает вас на 1% лучше."),
		list:     container.NewVBox(),
	}

	t.habits, err = storage.Load(dataPath)
	if err != nil {
		// Повреждённый файл откладываем в .bak, чтобы новые сохранения его не затёрли.
		if renameErr := os.Rename(dataPath, dataPath+".bak"); renameErr == nil {
			err = fmt.Errorf("%w\n\nСтарый файл сохранён как %s", err, dataPath+".bak")
		}
		t.habits = []models.Habit{}
		dialog.ShowError(err, window)
	}

	addButton := widget.NewButtonWithIcon("Добавить привычку", theme.ContentAddIcon(), t.showAddDialog)

	t.refresh()

	content := container.NewBorder(
		container.NewVBox(t.status, addButton, widget.NewSeparator()), // сверху
		nil, nil, nil,
		container.NewVScroll(t.list), // центр — список со скроллом
	)
	window.SetContent(content)
	window.ShowAndRun()
}
