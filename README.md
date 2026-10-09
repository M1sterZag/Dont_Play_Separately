# 🎮 Dont Play Separately (DPS) Backend

Backend REST-API для платформы — поиск тиммейтов для игр.

---

## 📋 Содержание

- [Технологии](#-технологии)
- [Архитектура](#-архитектура)
- [Структура проекта](#-структура-проекта)
- [Особенности](#-особенности)
- [Схема базы данных](#-схема-базы-данных)
- [API Эндпоинты](#-api-эндпоинты)
- [Быстрый старт](#-быстрый-старт)
- [Полезные команды](#-полезные-команды)

---

## 🛠 Технологии

| Компонент | Технология | Назначение |
|-----------|------------|------------|
| **Backend** | Go | Основной язык программирования |
| **Database** | PostgreSQL | Основное хранилище данных |
| **Migrations** | golang-migrate | Версионирование схемы БД ([migrations/](migrations)) |
| **Cache** | Redis | Кэширование данных, хранение refresh-токенов и кодов верификации |
| **Storage** | MinIO (S3) | Хранение аватаров пользователей и аватаров по умолчанию |
| **External API** | IGDB API | Получение данных об играх и игровых платформах |
| **Mail** | SMTP / MailHog | Отправка кодов подтверждения email |
| **Documentation** | Swagger/OpenAPI | Документация API ([docs/](docs)) |
| **Auth** | JWT | Аутентификация и авторизация |
| **Logging** | Uber Zap | Структурированное логирование |
| **Containerization** | Docker Compose | Локальная инфраструктура ([docker-compose.yaml](docker-compose.yaml)) |
| **CI** | GitHub Actions | Сборка, vet, тесты и линтер ([.github/workflows/ci.yml](.github/workflows/ci.yml)) |

---

## 🏗 Архитектура

Проект построен с разделением на слои:

```
internal/
├── core/                 # Инфраструктурное ядро
│   ├── cache/            # Redis-клиент
│   ├── config/           # Загрузка конфигурации из env
│   ├── domain/           # Общие доменные модели
│   ├── errors/           # Общие ошибки
│   ├── logger/           # Логгер (zap)
│   ├── mailer/           # SMTP-рассылка
│   ├── provider/igdb/    # Клиент IGDB API
│   ├── repository/       # Пул соединений PostgreSQL
│   ├── storage/          # S3-хранилище (MinIO)
│   └── transport/        # HTTP-сервер, middleware, роутер, хелперы request/response
└── features/             # Бизнес-фичи (каждая по слоям)
    ├── auth/             # Регистрация, вход, JWT, верификация email
    ├── games/            # Поиск и карточки игр
    ├── platforms/        # Игровые платформы
    ├── teams/            # Команды и участники
    └── users/            # Профили пользователей и аватары
```

Каждая фича делится на слои:

| Слой | Путь | Ответственность |
|------|------|-----------------|
| **transport** | `transport` | HTTP-хендлеры, DTO, роутинг |
| **service** | `service` | Бизнес-логика |
| **repository** | `repository` | Доступ к данным |
| **config** | `config.go` | Конфигурация фичи из env |

---

## ⚡ Особенности

### 🔐 JWT Аутентификация

| Действие | Описание |
|----------|----------|
| **Регистрация** | Создание аккаунта + отправка кода подтверждения email по SMTP |
| **Вход** | Пользователь получает `access_token` (1 час) и `refresh_token` (7 дней) |
| **Использование** | `access_token` передаётся в заголовке `Authorization: Bearer {token}` |
| **Обновление** | `refresh_token` (хранится в PostgreSQL хэшем) позволяет получить новый `access_token` |
| **Выход** | Отзыв сессии и занесение токена в чёрный список в Redis |
| **Верификация** | Код подтверждения хранится в Redis с TTL, ограничено число попыток и частота повторной отправки |

### 🎮 Работа с IGDB API

| Этап | Действие |
|------|----------|
| **1** | При отсутствии игры или платформы в PostgreSQL выполняется запрос к IGDB |
| **2** | Полученные данные сохраняются в PostgreSQL |
| **3** | Последующие запросы выполняются из PostgreSQL |
| **4** | Redis используется как кэш с настраиваемым TTL (`GAMES_CACHE_TTL`, `PLATFORMS_CACHE_TTL`) |
| **5** | Локальный поиск по названию ускорен trigram-индексами (`pg_trgm`) |

### 👤 Профили и аватары

- Аватар по умолчанию выбирается из набора [assets/avatars](assets/avatars) (по категориям игр) и загружается в MinIO командой `make s3-init`.
- Пользователь может обновлять профиль (никнейм, bio, аватар) и удалять аккаунт.

### 👥 Работа с командами

- Пользователь может создавать команды для конкретной игры и платформы.
- Каждая команда относится к **одной игре и одной платформе**.
- Пользователь может состоять в нескольких командах.
- Для команды можно указать требуемый рейтинг (`desired_rating`).
- Команда содержит ограниченное количество свободных мест (`slots_total`).
- Владелец команды имеет права управления её участниками и настройками.

---

## 🗄️ Схема базы данных

Все таблицы находятся в схеме `dps`. Миграции: [migrations/](migrations).

### 📊 Структура таблиц

#### 👤 Users (Пользователи)

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `UUID` (PK) | Уникальный идентификатор пользователя |
| `version` | `INTEGER` | Версия записи для оптимистичной блокировки |
| `email` | `VARCHAR(255)` (UNIQUE) | Email пользователя |
| `hashed_password` | `VARCHAR(255)` | Хэшированный пароль |
| `nickname` | `VARCHAR(40)` | Отображаемое имя пользователя |
| `bio` | `TEXT` | Информация о пользователе |
| `avatar_url` | `TEXT` | Ссылка на аватар пользователя |
| `created_at` | `TIMESTAMPTZ` | Дата и время регистрации |

---

#### 🔑 Refresh Sessions (Сессии refresh-токенов)

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `UUID` (PK) | Идентификатор сессии |
| `user_id` | `UUID` (FK → Users) | Владелец сессии |
| `refresh_token_hash` | `TEXT` | Хэш refresh-токена |
| `expires_at` | `TIMESTAMPTZ` | Время истечения |
| `revoked_at` | `TIMESTAMPTZ` | Время отзыва (NULL — активна) |
| `created_at` | `TIMESTAMPTZ` | Дата и время создания |

---

#### 🎮 Games (Игры)

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `BIGINT` (PK) | ID игры в IGDB |
| `title` | `TEXT` | Название игры |
| `icon_url` | `TEXT` | Ссылка на иконку игры |
| `slug` | `TEXT` (UNIQUE) | Slug из IGDB |
| `checksum` | `TEXT` | Контрольная сумма из IGDB |
| `updated_at` | `BIGINT` | Timestamp обновления в IGDB |
| `synced_at` | `TIMESTAMPTZ` | Время последней синхронизации |

---

#### 🖥️ Platforms (Игровые платформы)

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `BIGINT` (PK) | ID платформы в IGDB |
| `title` | `TEXT` | Название платформы |
| `abbreviation` | `TEXT` | Сокращённое название платформы |
| `slug` | `TEXT` (UNIQUE) | Slug из IGDB |
| `icon_url` | `TEXT` | Ссылка на иконку платформы |
| `checksum` | `TEXT` | Контрольная сумма из IGDB |
| `updated_at` | `BIGINT` | Timestamp обновления в IGDB |
| `synced_at` | `TIMESTAMPTZ` | Время последней синхронизации |

---

#### ❤️ User Platforms (Любимые платформы пользователей)

Таблица связывает пользователей с выбранными ими любимыми игровыми платформами.

| Поле | Тип | Описание |
|------|-----|----------|
| `user_id` | `UUID` (PK, FK → Users) | ID пользователя |
| `platform_id` | `BIGINT` (PK, FK → Platforms) | ID платформы |

> **ℹ️ Примечание:**
> Один пользователь может выбрать несколько любимых платформ.

---

#### 🎯 Teams (Команды)

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | `UUID` (PK) | Уникальный идентификатор команды |
| `version` | `INTEGER` | Версия записи для оптимистичной блокировки |
| `owner_id` | `UUID` (FK → Users) | Владелец команды |
| `game_id` | `BIGINT` (FK → Games) | Игра команды |
| `platform_id` | `BIGINT` (FK → Platforms) | Платформа команды |
| `title` | `VARCHAR(100)` | Название команды (1–100 символов) |
| `description` | `TEXT` | Описание команды (1–1000 символов) |
| `is_rating_required` | `BOOLEAN` | Будет ли рейтинговая игра |
| `desired_rating` | `TEXT` | Желаемый рейтинг или ранг (до 200 символов, обязателен при `is_rating_required = TRUE`) |
| `contact_link` | `TEXT` | Ссылка для связи (до 2000 символов) |
| `slots_total` | `INTEGER` | Общее количество мест в команде (> 0) |
| `created_at` | `TIMESTAMPTZ` | Дата и время создания команды |
| `is_active` | `BOOLEAN` | Активна ли команда |

---

#### 👥 Team Members (Участники команд)

Таблица связывает пользователей с командами.

| Поле | Тип | Описание |
|------|-----|----------|
| `team_id` | `UUID` (PK, FK → Teams) | ID команды |
| `user_id` | `UUID` (PK, FK → Users) | ID пользователя |

---

## 🔌 API Эндпоинты

Базовый префикс: `/api/v1` · Swagger UI: `/swagger`

### 🔐 Auth

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/auth/register` | Регистрация нового пользователя |
| `POST` | `/api/v1/auth/login` | Вход в систему |
| `POST` | `/api/v1/auth/refresh` | Обновление access-токена |
| `POST` | `/api/v1/auth/logout` | Выход из системы |
| `POST` | `/api/v1/auth/verify-email` | Подтверждение email кодом |
| `POST` | `/api/v1/auth/resend-verification-code` | Повторная отправка кода подтверждения |
| `PATCH` | `/api/v1/auth/change-password` 🔒 | Смена пароля |

---

### 👤 Users

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/users/profile/{user_id}` 🔒 | Получение профиля пользователя по ID |
| `PATCH` | `/api/v1/users/profile/me` 🔒 | Обновление своего профиля |
| `DELETE` | `/api/v1/users/profile/me` 🔒 | Удаление своего аккаунта |

---

### 🎯 Teams

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/teams` 🔒 | Создание новой команды |
| `GET` | `/api/v1/teams` | Получение списка команд с фильтрацией |
| `GET` | `/api/v1/teams/{team_id}` | Получение команды по ID |
| `PUT` | `/api/v1/teams/{team_id}` 🔒 | Редактирование команды |
| `DELETE` | `/api/v1/teams/{team_id}` 🔒 | Удаление команды |

---

### 👥 Team Members

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/teams/{team_id}/join` 🔒 | Присоединение к команде |
| `POST` | `/api/v1/teams/{team_id}/leave` 🔒 | Выход из команды |
| `DELETE` | `/api/v1/teams/{team_id}/members/{user_id}` 🔒 | Удаление участника из команды |

---

### 🎮 Games

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/games/search` | Поиск игр по названию |
| `GET` | `/api/v1/games/{id}` | Получение игры по ID |

---

### 🖥️ Platforms

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/platforms/search` | Поиск платформ по названию |
| `GET` | `/api/v1/platforms/{id}` | Получение платформы по ID |

> 🔒 — требуется `Authorization: Bearer {access_token}`

---

### 📊 Параметры фильтрации

#### GET `/api/v1/teams`

| Параметр | Тип | Описание |
|----------|-----|----------|
| `search` | `STRING` | Поиск по названию команды |
| `game_id` | `INTEGER` | Фильтр по ID игры в IGDB |
| `platform_id` | `INTEGER` | Фильтр по ID платформы в IGDB |
| `limit` | `INTEGER` | Количество записей на страницу (по умолчанию 20) |
| `offset` | `INTEGER` | Смещение для пагинации |

#### GET `/api/v1/games/search`

| Параметр | Тип | Описание |
|----------|-----|----------|
| `search` | `STRING` | Поиск по названию игры |
| `limit` | `INTEGER` | Количество записей (по умолчанию 10) |
| `offset` | `INTEGER` | Смещение для пагинации |

#### GET `/api/v1/platforms/search`

| Параметр | Тип | Описание |
|----------|-----|----------|
| `search` | `STRING` | Поиск по названию платформы |
| `limit` | `INTEGER` | Количество записей (по умолчанию 10) |
| `offset` | `INTEGER` | Смещение для пагинации |

---

## 🚀 Быстрый старт

### Предварительные требования

- [Go](https://go.dev/dl/) ≥ 1.26
- [Docker](https://docs.docker.com/get-docker/) + Docker Compose
- `openssl` (для генерации секретов)

### Запуск

1. **Создайте файл окружения** (Makefile и docker-compose читают `.env`):

   ```bash
   cp .env.example .env
   ```

   Заполните значения `POSTGRES_*`, `S3_*`, `JWT_SECRET`. Секрет можно сгенерировать:

   ```bash
   make generate-secret
   ```

2. **Поднимите инфраструктуру** (PostgreSQL, Redis, MinIO, MailHog + port-forwarders):

   ```bash
   make docker-up
   ```

3. **Примените миграции:**

   ```bash
   make migrate-up
   ```

4. **Инициализируйте S3-бакет** и загрузите аватары по умолчанию:

   ```bash
   make s3-init
   ```

5. **Запустите приложение:**

   ```bash
   make app-run
   ```

   API доступно на `http://localhost:5050`, Swagger UI — `http://localhost:5050/swagger`.

6. **(Опционально) перегенерируйте Swagger-документацию:**

   ```bash
   make swagger-gen
   ```

### Почта для разработки

Все письма попадают в MailHog: UI доступен на `http://localhost:8025` (SMTP — `localhost:1025`).

---

## 🧰 Полезные команды

| Команда | Описание |
|---------|----------|
| `make docker-up` | Запуск инфраструктуры (PostgreSQL, Redis, MinIO, MailHog) |
| `make docker-down` | Остановка всей инфраструктуры |
| `make volume-cleanup` | ⚠️ Удаление volume-данных PostgreSQL и MinIO |
| `make logs-cleanup` | ⚠️ Очистка файлов логов (`data/logs/`) |
| `make migrate-create name=<name>` | Создание новой миграции |
| `make migrate-up` | Применение миграций |
| `make migrate-down` | Откат миграций |
| `make port-forward-start` | Запуск port-forwarder'ов для локального доступа к сервисам |
| `make port-forward-stop` | Остановка port-forwarder'ов |
| `make s3-init` | Создание бакета и загрузка аватаров в MinIO |
| `make swagger-gen` | Генерация Swagger-документации |
| `make generate-secret` | Генерация случайного секретов и паролей для сервисов |
| `make app-run` | Запуск приложения в режиме разработки |

---

## 🚢 Деплой на сервер

Деплой автоматизирован через GitHub Actions: [deploy.yml](.github/workflows/deploy.yml).

### Как это работает

```text
push в main
     ↓
Job 1: сборка Docker-образа → публикация в GHCR (теги latest и sha-<коммит>)
     ↓
Job 2: rsync (compose, миграции, assets) → SSH на сервер
     ↓
Генерация .env из GitHub Secrets → docker compose up -d
```

Ключевые файлы:

| Файл | Назначение |
|------|------------|
| [Dockerfile](Dockerfile) | Многоэтапная сборка: компиляция в `golang:alpine`, минимальный runtime на `alpine` |
| [docker-compose.prod.yaml](docker-compose.prod.yaml) | Production-стек: app, PostgreSQL, Redis, MinIO, разовые миграции и init S3 |
| [.github/workflows/deploy.yml](.github/workflows/deploy.yml) | CI/CD деплой |

Переменные окружения на сервере создаются **автоматически**: workflow генерирует `.env` из GitHub Secrets при каждом деплое (в контейнерах сервисы доступны по именам docker-сети: `postgres`, `redis`, `minio:9000`).

### Подготовка (один раз)

1. Установить на сервере Docker Engine и Docker Compose plugin.
2. Создать ключ деплоя:

   ```bash
   ssh-keygen -t ed25519 -f dps_deploy_key -N ""
   # публичный ключ добавить в ~/.ssh/authorized_keys на сервере
   ```

3. На сервере: `mkdir -p /opt/dps` и дать пользователю запись в эту директорию.
4. В репозитории: **Settings → Secrets and variables → Actions** — добавить секреты из таблицы ниже.
5. Убедиться, что репозиторий и GHCR-пакет доступны (для приватного пакета сервер должен быть залогинен в GHCR либо пакет сделан public).

### Обязательные GitHub Secrets

| Секрет | Описание |
|--------|----------|
| `SSH_HOST` | IP или домен сервера |
| `SSH_USER` | Пользователь для SSH-деплоя |
| `SSH_PRIVATE_KEY` | Содержимое приватного ключа `dps_deploy_key` |
| `POSTGRES_PASSWORD` | Пароль БД |
| `JWT_SECRET` | Секрет для подписи JWT (`make generate-secret`) |
| `S3_ACCESS_KEY` | Логин MinIO |
| `S3_SECRET_KEY` | Пароль MinIO |
| `S3_PUBLIC_BASE_URL` | Публичный URL S3, например `http://<IP_сервера>:9000` |
| `IGDB_CLIENT_ID` | Client ID Twitch/IGDB |
| `IGDB_CLIENT_SECRET` | Client Secret Twitch/IGDB |
| `SMTP_HOST`, `SMTP_USER`, `SMTP_PASSWORD` | Данные почтового сервера для кодов верификации |
| `REDIS_PASSWORD` | Пароль Redis (`openssl rand -base64 24`) |

### Опциональные GitHub Secrets (есть дефолты)

| Секрет | Дефолт |
|--------|--------|
| `POSTGRES_USER` / `POSTGRES_DB` | `dps` |
| `S3_BUCKET` | `avatars` |
| `SMTP_PORT` | `587` |
| `SMTP_FROM` | `no-reply@dps.local` |
| `HTTP_ALLOWED_ORIGINS` | `*` |
| `APP_PORT` | `8080` |
| `TIME_ZONE` / `LOGGER_LEVEL` | `UTC` / `INFO` |

### Откат

Каждый образ публикуется с тегом `sha-<короткий хеш коммита>`. Откат:

```bash
ssh <user>@<host>
cd /opt/dps
sed -i 's|^APP_IMAGE=.*|APP_IMAGE=ghcr.io/<owner>/<repo>:sha-<хеш>|' .env
docker compose -f docker-compose.prod.yaml up -d
```
