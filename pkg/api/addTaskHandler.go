package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
	"github.com/IlyaOnyshchenko/go_final_project/pkg/nextdate"
)

// Обработчик для добавления задачи
func addTaskHandler(w http.ResponseWriter, req *http.Request) {
	// В task запишем результат десериализации
	var task db.Task
	// Десериализуем тело ответа в структуру task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "error JSON decoding while adding task"}, http.StatusBadRequest)
		return
	}
	// Проверяем, что поле Title не пустое
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "task title is empty for adding task"}, http.StatusBadRequest)
		return
	}
	// Проверяем на корректность полученное значение task.Date
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	// Присваеиваем id новой задаче
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	idOut := strconv.FormatInt(id, 10)
	// Возвращаем сообщение об успешно созданной задаче с id
	writeJson(w, map[string]string{
		"id": idOut,
	}, http.StatusOK)
}

// Функция проверки правильности введённой даты
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана, используем текущую дату
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	// Проверяем корректность формата даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("wrong date format, not like '20060102': ")
	}

	// Если дата в прошлом
	if afterNow(now, t) {
		if task.Repeat == "" {
			// Если правила повторения нет, используем сегодняшнюю дату
			task.Date = now.Format(DateFormat)
		} else {
			// Если есть правило повторения, вычисляем следующую дату
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("wrong input of repeat rule: %w", err)
			}
			task.Date = next
		}
	}

	return nil
}

// afterNow проверяет, что первая дата больше второй (игнорируя время)
func afterNow(date, now time.Time) bool {
	// Нормализуем даты, убирая время
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}
