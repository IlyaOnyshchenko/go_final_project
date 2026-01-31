package api

import (
	"net/http"
	"time"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
	"github.com/IlyaOnyshchenko/go_final_project/pkg/nextdate"
)

// Обработчик для редактирования задачи
func taskDoneHandler(w http.ResponseWriter, req *http.Request) {
	// Проверяем метод HTTP-запроса http.MethodPost()
	if req.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	// Получаем параметр search из URL
	id := req.URL.Query().Get("id")
	if id == "" {
		// Возвращаем результат в формате JSON
		writeJson(w, map[string]string{"error": "no id for task done"}, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	// Обрабатываем ошибку запроса к БД
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Если периодичность не указана, удаляем задачу
	if task.Repeat == "" {
		err := db.DeleteTask(id)
		// Обрабатываем ошибку запроса к БД
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJson(w, map[string]string{}, http.StatusOK)
		return
	}
	now := time.Now()
	nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
	// После вычисления следующей даты
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Обновляем дату выполнения в БД
	err = db.UpdateDate(nextDate, id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Возвращаем результат в формате JSON
	writeJson(w, map[string]string{}, http.StatusOK)
}
