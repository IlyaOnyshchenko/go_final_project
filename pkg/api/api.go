package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Init регистрирует все API обработчики
func Init() {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(taskDoneHandler))
}

// taskHandler обрабатывает запросы к /api/task в зависимости от HTTP-метода
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Структура для запроса аутентификации
// Описывает JSON‑формат входящего запроса на /api/signin
type SigninRequest struct {
	Password string `json:"password"`
}

// Структура для ответа аутентификации
// Определяет формат JSON‑ответа: либо токен, либо ошибка
type SigninResponse struct {
	Token string `json:"token,omitempty"` // JWT‑токен (при успешной аутентификации)
	Error string `json:"error,omitempty"` // Сообщение об ошибке (при неудаче)
}

// Секретный ключ для подписи JWT
var jwtKey = []byte("my_secret_key")

// Обработчик HTTP‑запроса на аутентификацию (/api/signin)
func signinHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем пароль из переменной окружения
	todoPassword := os.Getenv("TODO_PASSWORD")
	todoPassword = "12345" // Временное жёсткое задание пароля (для тестирования)

	// Если пароль не задан в окружении, возвращаем ошибку сервера
	if todoPassword == "" {
		http.Error(w, "Authentication not configured", http.StatusInternalServerError)
		return
	}

	// Читаем и декодируем JSON‑тело запроса в структуру SigninRequest
	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid request body")
		return
	}

	// Проверяем соответствие введённого пароля ожидаемому
	if req.Password != todoPassword {
		sendError(w, "Неверный пароль")
		return
	}

	// Генерируем JWT‑токен на основе пароля
	token, err := generateJWT(todoPassword)
	if err != nil {
		sendError(w, "Failed to generate token")
		return
	}

	// Отправляем успешный ответ с токеном в формате JSON
	json.NewEncoder(w).Encode(SigninResponse{Token: token})
}

// Вспомогательная функция для отправки JSON‑ответа с ошибкой
func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SigninResponse{Error: message})
}

// Функция генерации JWT‑токена
func generateJWT(password string) (string, error) {
	// Создаём токен с заявками (claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": hashPassword(password),               // Хэш пароля для проверки при валидации
		"exp":           time.Now().Add(8 * time.Hour).Unix(), // Время жизни токена
	})

	// Подписываем токен секретным ключом и получаем строку
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Функция валидации JWT‑токена
func validateJWT(tokenString string, expectedPassword string) bool {
	claims := jwt.MapClaims{}

	// Парсим токен и заполняем claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil // Предоставляем ключ для проверки подписи
	})

	// Если ошибка или токен невалиден — возвращаем false
	if err != nil || !token.Valid {
		return false
	}

	// Получаем хэш пароля из токена
	passwordHash, ok := claims["password_hash"].(string)
	if !ok {
		return false
	}

	// Сравниваем хэш из токена с хэшем текущего пароля
	return passwordHash == hashPassword(expectedPassword)
}

// Функция хеширования пароля
func hashPassword(password string) string {
	return password
}

// Middleware для проверки аутентификации
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем пароль из окружения
		todoPassword := os.Getenv("TODO_PASSWORD")

		// Если пароль не задан — пропускаем аутентификацию (режим без защиты)
		if todoPassword == "" {
			next(w, r)
			return
		}

		// Пытаемся получить куку token из запроса
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		token := cookie.Value // Значение токена из куки

		// Валидируем токен: проверка подписи и соответствие пароля
		if !validateJWT(token, todoPassword) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		// Если токен валиден — передаём запрос обработчику
		next(w, r)
	})
}
