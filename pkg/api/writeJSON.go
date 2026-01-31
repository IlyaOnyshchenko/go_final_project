package api

import (
	"encoding/json"
	"net/http"
)

// writeJson — вспомогательная функция для отправки JSON-ответа клиенту.
func writeJson(w http.ResponseWriter, data any, status int) error {
	// Устанавливаем заголовок Content-Type, указывающий, что ответ содержит JSON в кодировке UTF-8
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Отправляем указанный HTTP-статус-код в заголовке ответа
	w.WriteHeader(status)

	// Создаём JSON-энкодер, привязанный к ResponseWriter и кодируем данные в JSON
	return json.NewEncoder(w).Encode(data)
}
