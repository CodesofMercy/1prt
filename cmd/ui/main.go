package main

import (
	"1prc/internal/models"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var habits []models.Habit

// Функция для отметки привычки как выполненной
func markHabitAsDone(habitName string) {
	for i, habit := range habits {
		if habit.Name == habitName {
			habits[i].Completed = append(habits[i].Completed, time.Now())
			break
		}
	}
}

func habitExists(habitName string) bool {
	for _, habit := range habits {
		if strings.ToLower(habit.Name) == strings.ToLower(habitName) {
			return true
		}
	}
	return false
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("1% Habit Tracker")

	myWindow.Resize(fyne.NewSize(800, 600)) // Устанавливаем размер окна

	label := widget.NewLabel("Добро пожаловать в приложение 1%!")

	// Создаем контейнер для отображения привычек
	habitList := container.NewVBox()

	// Функция для обновления списка привычек на экране
	updateHabitList := func() {
		habitList.Objects = nil // Очищаем старый список
		for _, habit := range habits {
            habitLabel := widget.NewLabel(habit.Name)
        
            var markDoneButton *widget.Button
        
            // Проверяем, выполнена ли привычка
            if len(habit.Completed) > 0 {
                markDoneButton = widget.NewButton("Привычка выполнена", nil)
                markDoneButton.Disable() // Делаем кнопку неактивной, если привычка выполнена
            } else {
                markDoneButton = widget.NewButton("Отметить выполненной", func(habitName string) func() {
                    return func() {
                        markHabitAsDone(habitName)
                        label.SetText("Привычка выполнена: " + habitName)
                        markDoneButton.SetText("Привычка выполнена")
                        markDoneButton.Disable()
                    }
                }(habit.Name))
            }
        
            habitItem := container.NewHBox(habitLabel, markDoneButton)
            habitList.Add(habitItem)
        }
        habitList.Refresh() // Обновляем отображение списка
    }

	// Кнопка для добавления новой привычки
	button := widget.NewButton("Добавить привычку", func() {
		// Открытие диалога для добавления привычки
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("Введите название привычки")

		// Обработчик для нажатия Enter
		nameEntry.OnSubmitted = func(text string) {
			if text != "" {
				text = strings.TrimSpace(text)                                 // Убираем лишние пробелы
				capitalizedText := strings.ToUpper(string(text[0])) + text[1:] // Меняем только первую букву
				habit := models.Habit{
					Name:      capitalizedText,
					StartDate: time.Now(),
					Completed: []time.Time{},
				}
				habits = append(habits, habit)
				updateHabitList() // Обновляем список привычек
				label.SetText("Привычка добавлена: " + habit.Name)
				nameEntry.SetText("") // Очищаем поле после добавления
			}
		}

		dialog.ShowForm("Новая привычка", "Добавить", "Отмена",
			[]*widget.FormItem{
				widget.NewFormItem("Название", nameEntry),
			},
			func(b bool) {
				if b && nameEntry.Text != "" {
					capitalizedText := strings.ToUpper(string(nameEntry.Text[0])) + nameEntry.Text[1:]
					habit := models.Habit{
						Name:      capitalizedText,
						StartDate: time.Now(),
						Completed: []time.Time{},
					}
					habits = append(habits, habit)
					updateHabitList() // Обновляем список привычек
					label.SetText("Привычка добавлена: " + habit.Name)
					nameEntry.SetText("") // Очистка поля после добавления
				}
			}, myWindow)
	})

	// Контейнер для размещения всех элементов
	content := container.NewVBox(label, button, habitList)
	myWindow.SetContent(content)

	myWindow.ShowAndRun() // Запускаем приложение
}
