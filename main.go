// Основной пакет для запуска БД, обработчиков и сервера
package main

import (
	"log"

	"github.com/IlyaOnyshchenko/go_final_project/pkg/api"
	"github.com/IlyaOnyshchenko/go_final_project/pkg/db"
	"github.com/IlyaOnyshchenko/go_final_project/pkg/server"
)

func main() {
	// Запускаем и открываем соединение с БД
	err := db.Init("scheduler.db")
	// Обрабатываем ошибку
	if err != nil {
		log.Fatal(err)
	}
	// Запускаем обработчики
	api.Init()
	// Запускаем сервер
	server.ServerStart()
}
