# 1% Habit Tracker

Настольное приложение для отслеживания привычек по принципу «на 1% лучше каждый день».
Написано на Go с GUI на [Fyne](https://fyne.io).

**Версия: 0.1.0 (MVP)**

## Возможности

- Добавление привычек (с защитой от дубликатов)
- Отметка «выполнено сегодня» — на следующий день кнопка снова активна
- Подсчёт серии (streak): сколько дней подряд привычка выполняется
  (серия не сгорает, пока текущий день не закончился)
- Удаление привычки с подтверждением
- Данные сохраняются на диск и переживают перезапуск приложения

## Где хранятся данные

| ОС      | Путь                                        |
|---------|---------------------------------------------|
| Linux   | `~/.config/1prt/habits.json`                |
| macOS   | `~/Library/Application Support/1prt/habits.json` |
| Windows | `%AppData%\1prt\habits.json`                |

Формат — обычный JSON, файл можно бэкапить или переносить между машинами.

## Сборка и запуск

Требуется Go 1.23+ и компилятор C (Fyne использует cgo).

### Linux

```bash
sudo apt install gcc libgl1-mesa-dev xorg-dev   # Debian/Ubuntu
go build -o 1prt ./cmd/ui
./1prt
```

### Windows

Установите [Go](https://go.dev/dl/) и компилятор C (проще всего [MSYS2](https://www.msys2.org/) с mingw-w64, см. [инструкцию Fyne](https://docs.fyne.io/started/)), затем:

```powershell
go build -ldflags -H=windowsgui -o 1prt.exe ./cmd/ui
.\1prt.exe
```

### macOS

Нужен Xcode Command Line Tools (`xcode-select --install`), затем:

```bash
go build -o 1prt ./cmd/ui
./1prt
```

### Быстрый запуск без сборки

```bash
go run ./cmd/ui
```

## Тесты

```bash
go test ./...
```

## Структура проекта

```
cmd/ui/            — точка входа и GUI (Fyne)
internal/models/   — модель привычки: отметки по дням, подсчёт серий
internal/storage/  — сохранение/загрузка JSON с атомарной записью
```

## Планы после 0.1

- Редактирование привычки (имя, описание)
- История и календарь отметок
- Статистика и графики прогресса
- Напоминания / иконка в трее
- Сборка релизных бинарников через `fyne package` / GitHub Actions
