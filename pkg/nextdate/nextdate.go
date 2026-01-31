// Пакет nextdate создан отдельно для хранения функции вычисления следующей даты для задач с или без правил повторения
package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Функция NextDate вычисляет следующую дату, соответствующую правилу повторения.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим начальную дату из строки формата "YYYYMMDD"
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date")
	}

	// Проверяем, что правило повторения не пустое
	if repeat == "" || strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	// Разбиваем правило на части по пробелам
	rules := strings.Fields(repeat)
	if len(rules) == 0 {
		return "", fmt.Errorf("empty repeat rule")
	}

	// Обрабатываем разные типы правил повторения через switch
	switch rules[0] {
	// Правило "y": ежегодное повторение
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}
	// Правило "d": ежедневное повторение
	case "d":
		// Правило повторения дней должно сопровождаться числом после пробела
		if len(rules) < 2 {
			return "", fmt.Errorf("day interval required")
		}

		// Парсим интервал из второго элемента правила
		interval, err := strconv.Atoi(rules[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}
	// Правило "w": еженедельное повторение
	case "w":
		// Правило повторения дней недели должно сопровождаться числом после пробела
		if len(rules) < 2 {
			return "", fmt.Errorf("weekday list required")
		}

		// Помещаем в мапу количество допустимых дней недели
		target := map[int]bool{}
		for _, s := range strings.Split(rules[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("invalid weekday")
			}
			target[n] = true
		}

		// Начинаем поиск со следующего дня после `now`
		day := now.AddDate(0, 0, 1)
		// Устанавливаем лимит поиска — 1 год вперёд (чтобы избежать бесконечного цикла)
		limit := now.AddDate(1, 0, 0)

		for !day.After(limit) {
			// Проверяем, соответствует ли текущий день недели правилу
			if target[toRuleWeekday(day.Weekday())] {
				return day.Format("20060102"), nil
			}
			day = day.AddDate(0, 0, 1) // Переходим к следующему дню
		}

		return "", fmt.Errorf("no weekday match")
	// Правило "m": еженедельное повторение
	case "m":
		// Правило повторения месяцев должно сопровождаться числом(-ами) после пробела
		if len(rules) < 2 {
			return "", fmt.Errorf("month days required")
		}

		// Парсим дни месяца (например, "15,-1" → [15, -1])
		days := []int{}
		for _, s := range strings.Split(rules[1], ",") {
			n, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || n == 0 || n < -2 || n > 31 {
				return "", fmt.Errorf("invalid month day")
			}
			days = append(days, n)
		}

		// Парсим месяцы
		months := map[int]bool{}
		if len(rules) > 2 {
			for _, s := range strings.Split(rules[2], ",") {
				n, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("invalid month")
				}
				months[n] = true
			}
		}

		// Определяем начальную точку поиска
		start := date
		if now.After(start) {
			start = now
		}
		day := start.AddDate(0, 0, 1)

		// Ищем подходящую дату в пределах 1000 дней
		for i := 0; i < 1000; i++ {
			// Если указаны месяцы, проверяем, что текущий месяц входит в список
			if len(months) > 0 && !months[int(day.Month())] {
				day = day.AddDate(0, 0, 1)
				continue
			}

			// Получаем последний день текущего месяца
			lastDay := time.Date(day.Year(), day.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			// Проверяем каждый указанный день месяца
			for _, x := range days {
				want := x
				if x == -1 { // -1 означает последний день месяца
					want = lastDay
				}
				if x == -2 { // -2 означает предпоследний день месяца
					want = lastDay - 1
				}
				if day.Day() == want { // Если текущий день совпадает с нужным — нашли результат
					return day.Format("20060102"), nil
				}
			}

			day = day.AddDate(0, 0, 1) // Переходим к следующему дню
		}

		return "", fmt.Errorf("no monthly match")

	default:
		// Неизвестное правило повторения
		return "", fmt.Errorf("unknown rule")
	}
}

// Вспомогательная функция toRuleWeekday преобразует стандартный time.Weekday (0=воскресенье, ..., 6=суббота)
// в формат правила (1=понедельник, ..., 7=воскресенье).
func toRuleWeekday(wd time.Weekday) int {
	return int((wd+6)%7) + 1
}
