package api

import (
	"encoding/json"
	"net/http"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, req *http.Request) {
	// В task запишем результат десиреализации
	var task db.Task
	// Десериализуем тело ответа в структуру task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid JSON while updating task"}, http.StatusBadRequest)
		return
	}
	// Проверяем, что поле Title не пустое
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is empty"}, http.StatusBadRequest)
		return
	}
	// Проверяем, что поле ID не пустое
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "no id for task update"}, http.StatusBadRequest)
		return
	}
	// Проверяем на корректность полученное значение task.Date
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	// Присваеиваем id новой задаче
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Возвращаем пустое сообщение, означающее успешно обновленную задачу
	writeJson(w, map[string]string{}, http.StatusOK)
}
