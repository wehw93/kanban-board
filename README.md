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
- Получение логов задачи
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
1. Для начала нужно создать пользователя.
Введите данные в данный хендлер и нажмите execute

   ![image](https://github.com/user-attachments/assets/aacf6e4b-ed7a-4e81-a012-e5c7094984a2)
   
Ответ должен прийти такой

![image](https://github.com/user-attachments/assets/04438ffd-3b32-46c0-9ebb-94e8ea81fc94)

2. Далее нужно авторизоваться и получить JWT токен

![image](https://github.com/user-attachments/assets/32efc5fb-70af-4a59-9acb-141962b3b424)

Ответ должен быть такой. Копируем наш JWT токен.

![image](https://github.com/user-attachments/assets/0f229b60-a912-4b15-a468-0ce4e824fb54)

Теперь вставляем наш токен в authorize, первым слово пишем Bearer

![image](https://github.com/user-attachments/assets/6e6e5958-dea0-4e06-9e63-251ec60a7b3c)

Мы авторизованы.

3.Создаем проект

![image](https://github.com/user-attachments/assets/e7c30476-4930-4bbc-9066-4f2f79874395)

Ответ

![image](https://github.com/user-attachments/assets/ed07ab03-538d-4cdb-a929-65b44de6a032)

4.Далее нужно создать 3 колонки, или сколько вы хотите, но в данной методологии 3: todo, in progress, done

Вводим id проекта, который мы только что создали

![image](https://github.com/user-attachments/assets/e1ca1add-2701-48f6-9aa2-18f6f5369938)

![image](https://github.com/user-attachments/assets/1d221a24-d8be-4aee-8b4b-6c4c4a9cafcf)

![image](https://github.com/user-attachments/assets/e63af510-fe6d-410e-9bbc-e1ace48175f3)

5.Наконец создаем наши задачи

В поле id_column вводим id колонки, в которую мы хотим записать таск, этот айди вы можете найти в ответе, когда создавали колонки

![image](https://github.com/user-attachments/assets/b90db9ef-ea2c-422d-a9ae-50e84777ea9f)

Ответ

![image](https://github.com/user-attachments/assets/2db189e2-8f9c-4abe-8e1a-0ec4e24337fc)

Далее можем переместить нашуз задачу из одной колонки в другую

Здесь все параметры не обязательны, я хочу изменить статус задачи. Для этого нужно написать id колонки, в которую хочу переместить задачу.

![image](https://github.com/user-attachments/assets/fc79a1f2-8fe5-476d-b22f-e2fd3e2e73b8)

Ответ

![image](https://github.com/user-attachments/assets/d36d7686-a8ff-45db-b92e-937373f86f17)

Теперь статус нашей задачи изменился на "in_progress"

Таким образом можно перемещать задачи.

Каждое действие с задачей логируется и можно получить историю логов

![image](https://github.com/user-attachments/assets/d4c4e3ff-e8d3-4778-b3b3-76d35015bce6)

Ответ

![image](https://github.com/user-attachments/assets/f1660d2e-9833-494a-980d-aea9c17cfb50)










