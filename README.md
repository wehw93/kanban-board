# Kanban Board API: Your Ultimate Task Management System 📋🚀
## Kanban Board API — это мощный и удобный REST API сервис для управления задачами по методологии Kanban. С его помощью можно легко создавать проекты, колонки и задачи, перемещать их между статусами и отслеживать прогресс.
## ✨ Возможности
### Аутентификация
- Регистрация нового пользователя
- Вход в систему (получение JWT токена)
### Пользователи
- Просмотр информации о текущем пользователе
- Обновление данных пользователя
- Удаление пользователя
### Проекты
- Создание нового проекта
- Просмотр информации о проекте
- Обновление проекта
- Удаление проекта
- Получение списка всех проектов пользователя
### Колонки (Columns)
- Создание новой колонки в проекте
- Просмотр информации о колонке
- Обновление колонки
- Удаление колонки
### Задачи (Tasks)
- Создание новой задачи в колонке
- Просмотр информации о задаче
- Обновление задачи
- Удаление задачи
### Документация
- Доступ к Swagger документации API
## 📂 Структура проекта
```plaintext
.
├── board
├── cmd
│   ├── board
│   │   └── main.go
│   └── migrator
│       └── main.go
├── config
│   └── local.yaml
├── docker-compose.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── lib
│   │   ├── http
│   │   │   └── response
│   │   │       ├── column.go
│   │   │       ├── project.go
│   │   │       ├── response.go
│   │   │       └── user.go
│   │   ├── jwt
│   │   │   ├── helpers_jwt
│   │   │   │   └── helper_jwt.go
│   │   │   └── jwt.go
│   │   └── logger
│   │       └── sl
│   │           └── sl.go
│   ├── model
│   │   ├── column.go
│   │   ├── project.go
│   │   ├── task.go
│   │   ├── task_log.go
│   │   └── user.go
│   ├── service
│   │   ├── board
│   │   │   └── service.go
│   │   └── service.go
│   ├── storage
│   │   ├── board.go
│   │   ├── column.go
│   │   ├── postgresql
│   │   │   ├── column.go
│   │   │   ├── project.go
│   │   │   ├── storage.go
│   │   │   ├── task.go
│   │   │   ├── task_log.go
│   │   │   └── user.go
│   │   ├── storage.go
│   │   ├── task.go
│   │   ├── task_log.go
│   │   └── user.go
│   └── transport
│       └── http
│           ├── board.go
│           ├── column.go
│           ├── domain.go
│           ├── task.go
│           └── user.go
├── local.env
├── Makefile
├── migrations
│   ├── 1_create_users.down.sql
│   ├── 1_create_users.up.sql
│   ├── 2_create_projects.down.sql
│   ├── 2_create_projects.up.sql
│   ├── 3_create_column.down.sql
│   ├── 3_create_column.up.sql
│   ├── 4_create_tasks.down.sql
│   ├── 4_create_tasks.up.sql
│   ├── 5_create_logs.down.sql
│   └── 5_create_logs.up.sql
└── README.md
```
## 🚀 Установка и запуск
### Чтобы установить и запустить kanban доску, выполните следующие шаги:
1.Чтобы установить и запустить URLite, выполните следующие шаги:
```bash
https://github.com/wehw93/kanban-board.git
cd kanban-board
```
2.Установите зависимости:
```bash
go mod tidy
```
3.Настройте конфигурацию:
В папке config есть файл local.yaml. Убедитесь, что настройки корректны для вашего окружения:
```bash
env: "local" #prod

db:
  host: "localhost" #board_db
  port: "5433"
  name: "db_board"
  user: "board_user"
  password: "pwd123"
  sslmode: "disable"

http_server:
  address: "0.0.0.0:8080"
  timeout: "4s"
  idle_timeout: "60s"
```
4.Запустите докер контейнер для базы данных postgres
```bash
docker-compose up -d
```
5.Прогоните миграции
```bash
make run_migrations
```
6.Соберите приложение
```bash
make
```
7.Запустите приложение
```bash
./board
```
8.Откройте браузер и введите ссылку:
```bash
http://localhost:8080/swagger/index.html
```
9.Вы должны получить это:
![image](https://github.com/user-attachments/assets/30dca6ae-2f14-4ffe-b02b-498b7d27e6c7)

## 🛠 Использование
