# Workoutapp

Workoutapp — приложение для планирования и учёта тренировок. В нём можно создавать программы, планировать занятия и сохранять результаты тренировок. Приложение также показывает историю и статистику занятий.

## Основные возможности

- Регистрация и вход в аккаунт.
- Создание программ тренировок и тренировочных дней.
- Планирование занятий в календаре.
- Учёт выполненных упражнений и просмотр истории тренировок.
- Просмотр статистики тренировок.

## Технологический стек

- **Backend:** Go, `net/http`, pgx, JWT.
- **Frontend:** React, TypeScript, Vite, Tailwind CSS.
- **База данных:** PostgreSQL.
- **Инструменты:** Docker Compose, golang-migrate, Swagger.

## Структура проекта

```text
cmd/                  точка входа backend
internal/
  core/               общие доменные типы, middleware и инфраструктура
  features/           модули авторизации, программ, расписания и тренировок
migrations/           SQL-миграции базы данных
docs/swagger/         спецификация Swagger
frontend/
  src/
    components/       переиспользуемые компоненты интерфейса
    pages/            страницы приложения
    lib/              клиент API
  public/             статические ресурсы
```

Backend организован по слоям: HTTP-транспорт, сервисы и репозитории PostgreSQL. Frontend — React-приложение с отдельными страницами, компонентами и клиентом API.

## Локальный запуск

Необходимы Go **1.26.1**, Node.js с npm, Docker Compose и `make`.

1. Создайте файл `.env` на основе шаблона:

   ```sh
   cp .env.example .env
   ```

   Заполните значения в `.env`: `JWT_SECRET`, `HTTP_ADDR` (например, `:5050`), `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_TIMEOUT` (например, `5s`), `LOG_LEVEL` (например, `debug`) и `LOG_FOLDER` (например, `./out/logger`). Для локального запуска backend подключается к PostgreSQL на `localhost:5432`; порт по умолчанию — `5432`. Укажите также `UID` и `GID` текущего пользователя для Docker Compose.

2. Запустите PostgreSQL и примените миграции:

   ```sh
   make env-up
   make migrate-up
   ```

3. В корне проекта запустите backend:

   ```sh
   make app-run
   ```

   Эта цель задаёт `POSTGRES_HOST=localhost` и запускает приложение на адресе из `HTTP_ADDR`.

4. В другом терминале установите зависимости frontend и запустите Vite:

   ```sh
   cd frontend
   npm install
   npm run dev
   ```

   Приложение будет доступно по адресу [http://localhost:5173](http://localhost:5173). Запросы к backend frontend перенаправляет на `localhost:5050`. Swagger UI доступен по адресу [http://localhost:5050/swagger/index.html](http://localhost:5050/swagger/index.html).