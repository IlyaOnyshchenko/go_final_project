package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// Структура задачи в виде json
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Функция получения задачи по ID
func GetTask(id string) (*Task, error) {
	// Инициализируем idInt для корректного запроса в БД
	var idInt int64
	// Преобразуем в int64 переданный id
	idInt, err := strconv.ParseInt(id, 10, 64)
	// Обрабатываем ошибку
	if err != nil {
		return nil, fmt.Errorf("unknown id for task get: %s", id)
	}
	// Инициализируем в переменную task экземпляр структуры Task
	var task Task
	// В переменную query присваиваем текст запроса
	query := `SELECT * FROM scheduler WHERE id = ?`
	// Запрашиваем задачу по ID
	err = DB.QueryRow(query, idInt).Scan(&task.ID,
		&task.Date, &task.Title, &task.Comment, &task.Repeat,
	)
	// Обрабатываем ошибку запроса
	if err != nil {
		return nil, fmt.Errorf("request fault into DB while getting task: %w", err)
	}
	// Возвращаем экземпляр структуры с задачей
	return &task, nil
}

// Функция AddTask добавляет задачу в базу данных, возвращает id созданной задачи
func AddTask(task *Task) (int64, error) {
	// Определяем запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("request fault into DB while adding task:: %w", err)
	}

	// Получаем ID последней вставленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fault of getting last inserted id while adding task: %w", err)
	}
	return id, nil
}

// Функция UpdateTask редактирует задачу по переданной переменной task в виде структуры
func UpdateTask(task *Task) error {
	// Определяем запрос
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("request fault into DB while updating task: %w", err)
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("request error of rows affected while updating task: %w", err)
	}
	// Если count == 0, значит задача не была отредактирована
	if count == 0 {
		return fmt.Errorf("task not updated")
	}
	return nil
}

// Функция обновления даты повторяющейся задачи
func UpdateDate(next string, id string) error {
	// Инициализируем idInt для корректного запроса в БД
	var idInt int64
	// Преобразуем в int64 переданный id
	idInt, err := strconv.ParseInt(id, 10, 64)
	// Обрабатываем ошибку
	if err != nil {
		return fmt.Errorf("unknown id of task while updating: %s", id)
	}
	// Формируем запрос
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	// Исполняем запрос, заносим в res результат
	res, err := DB.Exec(query, next, idInt)
	if err != nil {
		return fmt.Errorf("request fault into DB while updating task: %w", err)
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("request error of rows affected while updating date: %w", err)
	}
	// Если count == 0, значит дата не была обновлена
	if count == 0 {
		return fmt.Errorf("task date not updated")
	}
	return nil
}

// Функция DeleteTask удаляет задачу по ID из БД
func DeleteTask(id string) error {
	// Инициализируем idInt для корректного запроса в БД
	var idInt int64
	// Преобразуем в int64 переданный id
	idInt, err := strconv.ParseInt(id, 10, 64)
	// Обрабатываем ошибку
	if err != nil {
		return fmt.Errorf("unknown id for task delete: %s", id)
	}
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, idInt)
	if err != nil {
		return fmt.Errorf("request fault into DB while deleting task: %w", err)
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("request error of rows affected while deleting task: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task not deleted")
	}
	return nil
}

// Функция Tasks возвращает список ближайших задач
// Также возвращает список задач по параметру search
func Tasks(limit int, search string) ([]*Task, error) {
	// Инициализируем пустой срез для хранения результатов
	res := make([]*Task, 0)

	// Объявляем переменные для работы с SQL-запросом
	var rows *sql.Rows
	var err error
	var query string
	var args []interface{}

	// Формируем SQL-запрос в зависимости от наличия поискового параметра
	if search == "" {
		// Если поиск не задан — получаем задачи без фильтрации, сортируем по дате и ограничиваем количество
		query = `SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		args = []interface{}{limit}
	} else {
		// Пытаемся интерпретировать поисковую строку как дату в формате "02.01.2006"
		parsed, parseErr := time.Parse("02.01.2006", search)
		if parseErr == nil {
			// Если строка успешно преобразована в дату — ищем задачи на конкретную дату
			searchDate := parsed.Format("20060102")
			query = `SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date, id LIMIT ?`
			args = []interface{}{searchDate, limit}
		} else {
			// Если строка не является датой — ищем по совпадению в заголовке или комментарии (частичное совпадение)
			query = `SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			searchLike := "%" + search + "%" // Добавляем % для поиска подстроки
			args = []interface{}{searchLike, searchLike, limit}
		}
	}

	// Выполняем SQL-запрос с переданными аргументами
	rows, err = DB.Query(query, args...)
	if err != nil {
		// Возвращаем ошибку, если запрос не удался
		return nil, fmt.Errorf("request fault into DB while getting task list: %w", err)
	}

	// Гарантируем закрытие набора строк после завершения функции
	defer rows.Close()

	// Проходим по всем полученным строкам
	for rows.Next() {
		// Создаём новый объект задачи
		task := &Task{}
		var id int64

		// Сканируем данные из текущей строки в поля задачи
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			// Возвращаем ошибку, если не удалось считать данные
			return nil, fmt.Errorf("request fault into DB while scaning rows: %w", err)
		}

		// Преобразуем числовой ID в строковый формат
		task.ID = strconv.FormatInt(id, 10)

		// Добавляем задачу в результирующий срез
		res = append(res, task)
	}

	// Проверяем, не возникла ли ошибка при чтении строк
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("request fault into DB after scaning rows: %w", err)
	}

	// Возвращаем результат (срез задач) и отсутствие ошибки
	return res, nil
}
