# 📚 Инструкция по использованию Auth Service

Полная документация по системе авторизации с поддержкой ролей и прав доступа.

---

## 📋 Содержание

1. [Описание проекта](#описание-проекта)
2. [Технологии](#технологии)
3. [Установка и запуск](#установка-и-запуск)
4. [Структура проекта](#структура-проекта)
5. [API Документация](#api-документация)
6. [Frontend](#frontend)
7. [База данных](#база-данных)
8. [Примеры использования](#примеры-использования)
9. [Безопасность](#безопасность)

---

## 🎯 Описание проекта

Auth Service — это полнофункциональная система авторизации и аутентификации, построенная на Go (Gin) и Nuxt 4. Система поддерживает:

- ✅ Регистрацию и авторизацию пользователей
- ✅ JWT токены с автоматическим обновлением (access token на 15 минут, refresh token на 7 дней)
- ✅ Систему ролей (admin, user, auditor)
- ✅ Гибкую систему прав доступа (read, edit, delete, create)
- ✅ Современный веб-интерфейс на Nuxt UI
- ✅ Личный кабинет пользователя

---

## 🛠 Технологии

### Backend
- **Go 1.21+** — основной язык программирования
- **Gin** — HTTP веб-фреймворк
- **PostgreSQL** — база данных
- **JWT (golang-jwt/v5)** — токены авторизации
- **sqlx** — работа с БД
- **bcrypt** — хеширование паролей

### Frontend
- **Nuxt 4** — Vue.js фреймворк
- **Nuxt UI** — компонентная библиотека
- **TypeScript** — типизация
- **Tailwind CSS** — стилизация

---

## 🚀 Установка и запуск

### Требования

- Docker и Docker Compose
- Go 1.21 или выше
- Node.js 18+ и pnpm

### Шаг 1: Запуск базы данных

```bash
docker-compose up -d
```

База данных будет доступна на порту `2020`.

### Шаг 2: Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```env
DB_DSN=postgres://user:password@localhost:2020/authdb?sslmode=disable
JWT_SECRET=myjwtsecret
SUPERADMIN_CODE=L@b47850811
```

### Шаг 3: Запуск Backend сервера

```bash
# Windows
set DB_DSN=postgres://user:password@localhost:2020/authdb?sslmode=disable
set JWT_SECRET=myjwtsecret
set SUPERADMIN_CODE=L@b47850811
go run .

# Linux/Mac
export DB_DSN=postgres://user:password@localhost:2020/authdb?sslmode=disable
export JWT_SECRET=myjwtsecret
export SUPERADMIN_CODE=L@b47850811
go run .
```

Сервер запустится на порту `8999`.

### Шаг 4: Запуск Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

Frontend будет доступен на `http://localhost:3000`.

---

## 📁 Структура проекта

```
auth-service-main/
├── main.go              # Точка входа, роутинг
├── handlers.go          # Обработчики HTTP запросов
├── middleware.go        # Middleware для авторизации
├── models.go            # Модели данных
├── db.go                # Инициализация БД и миграции
├── docker-compose.yml   # Конфигурация Docker
├── go.mod               # Зависимости Go
├── .env                 # Переменные окружения
│
└── frontend/
    ├── app/
    │   ├── pages/       # Страницы приложения
    │   │   ├── index.vue      # Главная страница
    │   │   ├── login.vue      # Страница входа
    │   │   ├── register.vue   # Страница регистрации
    │   │   └── me.vue         # Личный кабинет
    │   ├── components/  # Компоненты Vue
    │   │   └── AppHeader.vue  # Шапка сайта
    │   ├── composables/ # Композаблы
    │   │   └── useApi.ts      # API клиент
    │   └── assets/      # Статические файлы
    │       └── css/
    │           └── main.css    # Стили
    └── nuxt.config.ts   # Конфигурация Nuxt
```

---

## 📡 API Документация

### Базовый URL

```
http://localhost:8999
```

### Публичные эндпоинты

#### 1. Регистрация пользователя

```http
POST /register
Content-Type: application/json

{
  "username": "user1",
  "password": "password123",
  "role": "user"  // опционально, по умолчанию "user"
}
```

**Ответ:**
```json
{
  "message": "User registered"
}
```

#### 2. Авторизация

```http
POST /login
Content-Type: application/json

{
  "username": "user1",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Время жизни токенов:**
- `access_token` — 15 минут
- `refresh_token` — 7 дней

#### 3. Обновление Access Token

```http
POST /refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Защищенные эндпоинты

Все защищенные эндпоинты требуют заголовок:

```http
Authorization: Bearer <access_token>
```

#### 1. Получение информации о текущем пользователе

```http
GET /protected/me
Authorization: Bearer <access_token>
```

**Ответ:**
```json
{
  "username": "alisher",
  "role": "admin",
  "permissions": ["read", "edit", "delete", "create"]
}
```

#### 2. Эндпоинты по ролям

```http
GET /protected/user   # Только для роли "user"
GET /protected/admin  # Только для роли "admin"
```

#### 3. Эндпоинты по правам доступа

```http
GET /protected/read   # Требует право "read"
GET /protected/edit   # Требует право "edit"
GET /protected/delete # Требует право "delete"
```

### Управление ролями и правами

#### Получить список всех ролей и их прав

```http
GET /roles
```

**Ответ:**
```json
[
  {
    "role": "admin",
    "description": "Full access to everything",
    "permissions": ["read", "edit", "delete", "create"]
  },
  {
    "role": "user",
    "description": "Standard user",
    "permissions": ["read", "create"]
  },
  {
    "role": "auditor",
    "description": "Read-only access",
    "permissions": ["read"]
  }
]
```

#### Изменить права ролей

```http
POST /roles/update
Content-Type: application/json

{
  "superadmin_code": "L@b47850811",
  "updates": [
    {
      "role": "user",
      "permissions": ["read"]
    },
    {
      "role": "admin",
      "permissions": ["read", "edit", "delete", "create"]
    }
  ]
}
```

---

## 🎨 Frontend

### Страницы

1. **Главная страница (`/`)** — доступна всем пользователям
   - Для неавторизованных: информация о сервисе и кнопки входа/регистрации
   - Для авторизованных: приветствие и информация о пользователе

2. **Вход (`/login`)** — форма авторизации
   - Автоматическое сохранение токенов в cookies
   - Автоматическое обновление токенов при истечении

3. **Регистрация (`/register`)** — форма регистрации
   - Выбор роли при регистрации
   - Автоматическое перенаправление на страницу входа после успешной регистрации

4. **Личный кабинет (`/me`)** — информация о пользователе
   - Username и роль
   - Список прав доступа с визуальными индикаторами

### Компоненты

- **AppHeader** — шапка сайта с информацией об авторизации
- **useApi** — composable для работы с API
  - Автоматическое добавление токена в заголовки
  - Автоматическое обновление токена при 401 ошибке
  - Хранение токенов в cookies

### Цветовая схема

- **Основной цвет:** Фиолетовый (violet)
- **Поддержка темной темы:** Да
- **Адаптивный дизайн:** Да

---

## 🗄 База данных

### Структура таблиц

#### Таблица `users`

| Поле     | Тип    | Описание                             |
| -------- | ------ | ------------------------------------ |
| id       | SERIAL | Уникальный ID                        |
| username | TEXT   | Имя пользователя (уникальное)        |
| password | TEXT   | Хеш пароля (bcrypt)                  |
| role     | TEXT   | Название роли                        |

#### Таблица `roles`

| Поле        | Тип    | Описание      |
| ----------- | ------ | ------------- |
| id          | SERIAL | Уникальный ID |
| name        | TEXT   | Название роли |
| description | TEXT   | Описание роли |

#### Таблица `permissions`

| Поле        | Тип    | Описание                             |
| ----------- | ------ | ------------------------------------ |
| id          | SERIAL | Уникальный ID                        |
| name        | TEXT   | Название права                        |
| description | TEXT   | Описание действия                    |

#### Таблица `role_permissions`

| Поле          | Тип | Описание        |
| ------------- | --- | --------------- |
| role_id       | INT | Ссылка на роль  |
| permission_id | INT | Ссылка на право |

### Роли по умолчанию

При первом запуске автоматически создаются:

1. **admin** — полный доступ ко всем правам
   - Права: `read`, `edit`, `delete`, `create`

2. **user** — стандартный пользователь
   - Права: `read`, `create`

3. **auditor** — только чтение
   - Права: `read`

---

## 💡 Примеры использования

### Пример 1: Регистрация и авторизация

```bash
# Регистрация
curl -X POST http://localhost:8999/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123","role":"user"}'

# Авторизация
curl -X POST http://localhost:8999/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'

# Сохраните access_token из ответа
```

### Пример 2: Получение информации о пользователе

```bash
curl -X GET http://localhost:8999/protected/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Пример 3: Обновление токена

```bash
curl -X POST http://localhost:8999/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### Пример 4: Использование через Frontend

1. Откройте `http://localhost:3000`
2. Нажмите "Зарегистрироваться"
3. Заполните форму регистрации
4. После регистрации войдите в систему
5. Перейдите в личный кабинет для просмотра прав доступа

---

## 🔐 Безопасность

### Рекомендации

1. **JWT_SECRET** — используйте сложный секретный ключ в продакшене
2. **SUPERADMIN_CODE** — храните в безопасном месте, не коммитьте в репозиторий
3. **HTTPS** — используйте HTTPS в продакшене
4. **CORS** — настройте правильные домены для CORS
5. **Пароли** — используйте сложные пароли, система хеширует их через bcrypt

### Токены

- Access токены имеют короткое время жизни (15 минут) для безопасности
- Refresh токены используются для обновления access токенов
- Токены содержат информацию о пользователе и роли
- При истечении access токена фронтенд автоматически обновляет его через refresh токен

---

## 🐛 Решение проблем

### Проблема: База данных не подключается

**Решение:**
- Убедитесь, что Docker контейнер запущен: `docker ps`
- Проверьте переменную `DB_DSN` в `.env`
- Проверьте логи: `docker logs auth_postgres`

### Проблема: Токены не обновляются

**Решение:**
- Проверьте, что refresh токен сохранен в cookies
- Проверьте консоль браузера на ошибки
- Убедитесь, что сервер запущен и доступен

### Проблема: Права доступа не отображаются

**Решение:**
- Перезапустите backend сервер после изменений
- Проверьте, что права назначены роли в БД
- Проверьте логи сервера на ошибки

### Проблема: CORS ошибки

**Решение:**
- Убедитесь, что frontend запущен на `http://localhost:3000`
- Проверьте настройки CORS в `main.go`
- Добавьте нужный домен в `AllowOrigins`

---

## 📝 Дополнительная информация

### Тестовые пользователи

После запуска можно создать тестовых пользователей:

```bash
# Администратор
curl -X POST http://localhost:8999/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","role":"admin"}'

# Обычный пользователь
curl -X POST http://localhost:8999/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123","role":"user"}'
```

### Логирование

Backend логирует:
- Ошибки подключения к БД
- Ошибки получения прав доступа
- Запуск сервера

### Производительность

- Access токены проверяются на каждом защищенном запросе
- Права доступа кэшируются в JWT токене
- База данных использует индексы для быстрого поиска

---

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи сервера
2. Проверьте консоль браузера (F12)
3. Убедитесь, что все сервисы запущены
4. Проверьте переменные окружения

---

## 📄 Лицензия

Этот проект создан для образовательных целей.

---

**Версия:** 1.0.0  
**Последнее обновление:** 2025

