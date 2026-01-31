package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/nextdate"
)

// Константа для формата даты 20060102
const DateFormat = "20060102"

// NextDateHandler обрабатывает GET-запросы к /api/nextdate и возвращает следующую дату в формате 20060102 или текст ошибки
func NextDateHandler(w http.ResponseWriter, req *http.Request) {
	// Проверяем, что это GET-запрос
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из URL
	nowParam := req.FormValue("now")
	dateParam := req.FormValue("date")
	repeatParam := req.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateParam == "" {
		http.Error(w, "no date in input", http.StatusBadRequest)
		return
	}

	// Определяем время now
	var now time.Time
	if nowParam == "" {
		// Если параметр now не указан, используем текущую дату
		now = time.Now()
	} else {
		// Парсим переданную дату now
		parsedNow, err := time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("now is not parsed: %s", nowParam), http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	// Вызываем функцию NextDate из пакета nextdate
	nextDate, err := nextdate.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат в формате 20060102
	// Если nextDate пустой, возвращаем ошибку
	if nextDate == "" {
		http.Error(w, "unknown repeat rule", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}
