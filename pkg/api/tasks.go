package api

import (
	"net/http"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Обработчик для получения списка ближайших задач
func tasksHandler(w http.ResponseWriter, req *http.Request) {
	// Проверяем метод HTTP-запроса http.MethodGet()
	if req.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	// Получаем параметр search из URL
	search := req.URL.Query().Get("search")
	// В tasks помещаем результат запроса к БД
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Возвращаем результат в формате JSON
	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
