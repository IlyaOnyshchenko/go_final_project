package api

import (
	"net/http"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
)

// type Task struct {
// 	Task *db.Task `json:"task"`
// }
// Обработчик для получения задачи
func getTaskHandler(w http.ResponseWriter, req *http.Request) {
	// В task запишем результат запроса задачи по ID
	var task *db.Task
	// Получаем параметр search из URL
	id := req.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "no id for task get"}, http.StatusBadRequest)
		return
	}
	// В tasks помещаем результат запроса к БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Возвращаем результат в формате JSON
	writeJson(w, task, http.StatusOK)
}
