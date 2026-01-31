package api

import (
	"net/http"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
)

// Обработчик для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	// Получаем параметр search из URL
	id := req.URL.Query().Get("id")
	if id == "" {
		// Возвращаем результат в формате JSON
		writeJson(w, map[string]string{"error": "no id for task delete"}, http.StatusBadRequest)
		return
	}
	err := db.DeleteTask(id)
	// Обрабатываем ошибку запроса к БД
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}
