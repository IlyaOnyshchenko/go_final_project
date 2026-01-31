// Пакет server служит для запуска сервера с переменной окружения TODO_PORT, заданной программно, либо методом получения из переменной окружения
package server

import (
	"log"
	"net/http"
	"os"
)

// Функция запуска сервера
func ServerStart() {
	// Получаем значение переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")

	// Если переменная не задана — используем порт 7540 по умолчанию
	if port == "" {
		port = "7540"
		log.Println("TODO_PORT not set. Using default port:", port)
	} else {
		log.Println("Using port from TODO_PORT:", port)
	}
	// Указываем директорию, которую хотим сделать доступной для сервера
	fs := http.FileServer(http.Dir("./web"))

	// Регистрируем обработчик для корня хранилища файл-сервера "/"
	http.Handle("/", fs)

	// Запускаем сервер с использованием переменной окружения port
	log.Println("Starting server on :", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
