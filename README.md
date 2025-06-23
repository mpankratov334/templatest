### README.md для проекта на Go (DDD-стиль)

---

#### **Обзор проекта**
Проект представляет собой REST API для управления задачами (tasks), реализованный с использованием:
- **Domain-Driven Design (DDD)** для структурирования кода
- **In-memory хранилища** (с ограничениями по размеру)
- Веб-фреймворка **Fiber**
- Валидации запросов через кастомный валидатор
- Логирования через **Zap**
- Конфигурации через **.env-файл**

---

#### **Структура проекта**
```bash
.
├── cmd
│   └── main.go          # Точка входа
├── internal
│   ├── api
│   │   ├── api.go               # Роутеры
│   │   ├── middleware
│   │   │   └── authorization.go # Заглушка аутентификации
│   │   └── dto
│   │       └── http_responses.go # HTTP-ответы
│   ├── config
│   │   └── config.go    # Конфигурация приложения
│   ├── logger
│   │   └── logger.go    # Инициализация логгера (Zap)
│   ├── repo
│   │   ├── entity.go    # Сущности данных
│   │   └── repo.go      # In-memory репозиторий
│   └── service
│       ├── entity.go    # DTO для запросов
│       └── service.go   # Бизнес-логика
├── pkg
│   └── validator        # Кастомный валидатор (не показан)
└── local.env            # Файл конфигурации
```

---

#### **Зависимости**
Убедитесь, что установлены:
- [Go 1.18+](https://go.dev/dl/)
- Зависимости из `go.mod`:
  ```bash
  go get github.com/gofiber/fiber/v2
  go get github.com/joho/godotenv
  go get github.com/kelseyhightower/envconfig
  go get go.uber.org/zap
  ```

---

#### **Конфигурация (local.env)**
Создайте файл `local.env` в корне проекта:
```ini
# Настройки сервера
PORT=8080
REQUEST_TIMEOUT=30s

# Логирование
LOG_LEVEL=info

# In-memory хранилище
MAX_ITEMS=10000     # Макс. количество задач
MAX_ITEM_SIZE=1024  # Макс. размер задачи (байты)

# Аутентификация
API_TOKEN=secret    # Токен для заглушки аутентификации
```

---

#### **Запуск проекта**
1. Клонируйте репозиторий
2. Установите зависимости
3. Запустите приложение:
```bash
go run cmd/main.go
```
Сервер запустится на порту `8080` (настраивается в `local.env`).

---

#### **API Endpoints**
Все запросы требуют заголовок `Authorization: Bearer <API_TOKEN>` (токен из `local.env`).

| Метод  | Путь           | Действие               | Пример тела запроса (JSON)       |
|--------|----------------|------------------------|----------------------------------|
| POST   | /v1/tasks      | Создать задачу         | `{"title": "Task 1", "data": "..."}` |
| GET    | /v1/tasks      | Получить все задачи    | -                                |
| GET    | /v1/tasks/{id} | Получить задачу по ID  | -                                |
| PUT    | /v1/tasks/{id} | Обновить задачу        | `{"status": "done"}`             |
| DELETE | /v1/tasks/{id} | Удалить задачу         | -                                |

---

#### **Примеры запросов**
1. **Создание задачи**:
```bash
curl -X POST http://localhost:8080/v1/tasks \
  -H "Authorization: Bearer secret" \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "data": "Study DDD patterns"}'
```

2. **Получение всех задач**:
```bash
curl -H "Authorization: Bearer secret" http://localhost:8080/v1/tasks
```

---

#### **Валидация**
Запросы валидируются по правилам:
- `POST /tasks`:
  ```go
  type PostRequest struct {
    Title  string `json:"title" validate:"required"` // Обязательное поле
    Data   string `json:"data"`
    Status string `json:"status"`
  }
  ```
- Параметры пути (`/:id`):
  ```go
  type RequestWithId struct {
    ID string `validate:"required,intString,min=1"` // Число > 0
  }
  ```

При ошибке возвращается:
```json
{
  "status": "error",
  "error": {
    "code": "FIELD_INCORRECT",
    "desc": "Validation failed"
  }
}
```

---

#### **Ограничения хранилища**
- Максимальное количество задач: **10,000** (`MAX_ITEMS`)
- Максимальный размер одной задачи: **1 КБ** (`MAX_ITEM_SIZE`)
- При переполнении: `503 Service Unavailable`

---

#### **Логирование**
Настраивается через `LOG_LEVEL` в `local.env`:
- `debug` - детальные логи
- `info` - стандартный уровень
- `error` - только ошибки

Пример лога:
```json
{"timestamp":"2025-06-22T15:04:05.999Z","message":"Starting server on :8080"}
```

---

#### **Типовой ответ API**
Успех:
```json
{
  "status": "success",
  "data": {
    "task_id": 42
  }
}
```

Ошибка:
```json
{
  "status": "error",
  "error": {
    "code": "NOT_FOUND",
    "desc": "Task not found"
  }
}
```

---

#### **Дополнительно**
- Аутентификация: текущая реализация (`authorization.go`) является **заглушкой**
- In-memory данные **сбрасываются при перезапуске** сервера