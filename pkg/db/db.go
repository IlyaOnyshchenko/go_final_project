// В пакете db описаны функции запуска базы данных, а также методы работы с таблицой:
// получение списка ближайших задач, получения по id, удаления, добавления, обновления даты
package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Константы с строками запросов для создания таблицы, индекса даты
const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(32) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(32) NOT NULL DEFAULT ""
	)`
	addDI = `CREATE INDEX IF NOT EXISTS date_index ON scheduler (date)`
)

// Экземпляр экспортируемой переменной БД
var DB *sql.DB

// Функция запуска БД
func Init(dbFile string) error {
	// Открываем соединение с БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error DB opening: %w", err)
	}

	// Сохраняем соединение в глобальную переменную
	DB = db

	// Проверяем подключение
	if err := DB.Ping(); err != nil {
		DB.Close()
		return fmt.Errorf("error DB connection: %w", err)
	}

	// Создаём таблицу
	_, err = DB.Exec(schema)
	if err != nil {
		DB.Close()
		return fmt.Errorf("error table creating: %w", err)
	}

	// Создаём индекс
	_, err = DB.Exec(addDI)
	if err != nil {
		DB.Close()
		return fmt.Errorf("error index creating: %w", err)
	}

	return nil
}
