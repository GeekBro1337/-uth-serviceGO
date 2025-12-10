## 🧩 1. Логика связей между таблицами

### Основная схема

```
┌───────────┐        ┌────────────┐        ┌──────────────────┐
│  users    │        │   roles    │        │   permissions    │
│-----------│        │------------│        │------------------│
│ id        │  ◀──┐  │ id         │        │ id               │
│ username  │     │  │ name       │        │ name             │
│ password  │     │  │ description│        │ description      │
│ role      │──┐  │  └────────────┘        └──────────────────┘
└───────────┘  │  │
               │  │
               │  ▼
          ┌────────────────────┐
          │ role_permissions   │
          │--------------------│
          │ role_id            │──► roles.id
          │ permission_id      │──► permissions.id
          └────────────────────┘
```

### 🔗 Объяснение связей

- **users.role** → текстовое имя роли (например, `"admin"` или `"user"`)  
   Используется для простоты авторизации.
- **roles** → справочник всех ролей в системе.
- **permissions** → справочник всех возможных действий (например: `"read"`, `"edit"`, `"delete"`, `"create"`).
- **role_permissions** → связывает роли с разрешениями («админ может всё», «юзер — только читать»).

### Пример данных

| roles | permissions | role_permissions |
| ----- | ----------- | ---------------- |
| admin | read        | admin → read     |
| admin | edit        | admin → edit     |
| user  | read        | user → read      |
| user  | create      | user → create    |

---

````markdown
# Auth-Service (Go + Gin + PostgreSQL)

Многоуровневый API-сервис для регистрации, авторизации и управления ролями и правами пользователей.  
Использует JWT для токенов и PostgreSQL для хранения пользователей, ролей и разрешений.

---

## 🚀 Запуск

1. Подними базу данных PostgreSQL через Docker:

```bash
docker-compose up -d
```
````

2. Запусти сервис:

```bash
go run .
```

После старта сервер слушает порт `8999`.

---

## 🧩 Структура базы данных

### Таблица `users`

| Поле     | Тип    | Описание                             |
| -------- | ------ | ------------------------------------ |
| id       | SERIAL | Уникальный ID                        |
| username | TEXT   | Имя пользователя                     |
| password | TEXT   | Хеш пароля                           |
| role     | TEXT   | Название роли (`user`, `admin`, ...) |

### Таблица `roles`

| Поле        | Тип    | Описание      |
| ----------- | ------ | ------------- |
| id          | SERIAL | Уникальный ID |
| name        | TEXT   | Название роли |
| description | TEXT   | Описание роли |

### Таблица `permissions`

| Поле        | Тип    | Описание                             |
| ----------- | ------ | ------------------------------------ |
| id          | SERIAL | Уникальный ID                        |
| name        | TEXT   | Название права (`read`, `edit`, ...) |
| description | TEXT   | Описание действия                    |

### Таблица `role_permissions`

| Поле          | Тип | Описание        |
| ------------- | --- | --------------- |
| role_id       | INT | Ссылка на роль  |
| permission_id | INT | Ссылка на право |

---

## 🔑 Как связаны данные

- Каждый пользователь (`users.role`) имеет роль из таблицы `roles`
- Каждая роль имеет набор разрешений (`role_permissions`)
- Разрешения описаны в таблице `permissions`
- Проверка доступа происходит через middleware `PermissionMiddleware`

---

## 🧠 Пример логики

- Роль `admin` имеет все права (`read`, `edit`, `delete`, `create`)
- Роль `user` имеет только `read` и `create`
- Роль `auditor` имеет только `read`

---

## 🧾 API Эндпоинты

### 🧍 Регистрация

`POST /register`

```json
{
  "username": "user1",
  "password": "password123",
  "role": "user"
}
```

Ответ:

```json
{ "message": "User registered" }
```

---

### 🔑 Логин

`POST /login`

```json
{
  "username": "user1",
  "password": "password123"
}
```

Ответ:

```json
{ "token": "<JWT-токен>" }
```

---

### 🔒 Защищённые маршруты

Все требуют заголовок:

```
Authorization: Bearer <JWT-токен>
```

| Метод | Путь               | Описание                    |
| ----- | ------------------ | --------------------------- |
| `GET` | `/protected/user`  | Только для роли `user`      |
| `GET` | `/protected/admin` | Только для роли `admin`     |
| `GET` | `/protected/me`    | Возвращает данные из JWT    |
| `GET` | `/protected/read`  | Проверяет разрешение `read` |
| `GET` | `/protected/edit`  | Проверяет разрешение `edit` |

---

### ⚙️ Управление ролями и правами (через SUPERADMIN_CODE)

#### 🔹 Получить список всех ролей и их прав

`GET /roles`

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
  }
]
```

#### 🔹 Изменить права ролей

`POST /roles/update`

Тело запроса:

```json
{
  "superadmin_code": "MySecretAdminCode",
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

## 🧰 Middleware

| Название               | Назначение                                               |
| ---------------------- | -------------------------------------------------------- |
| `AuthMiddleware`       | Проверяет JWT, извлекает `username` и `role`             |
| `RoleMiddleware`       | Проверяет, что роль пользователя соответствует требуемой |
| `PermissionMiddleware` | Проверяет наличие конкретного разрешения у роли          |

---

## 🧾 Технологии

- **Go + Gin** — HTTP фреймворк
- **PostgreSQL + sqlx** — база данных
- **JWT (golang-jwt/v5)** — аутентификация
- **Docker** — база и локальный запуск

---

## 🔐 Переменные окружения (.env)

```env
DB_DSN=postgres://user:password@localhost:2020/authdb?sslmode=disable
JWT_SECRET=myjwtsecret
SUPERADMIN_CODE=L@b47850811
```

---

## ⚙️ Возможности для расширения

- Refresh Tokens
- Логирование действий
- Создание/удаление пользователей через API
- UI-панель управления ролями и правами

---

## 📊 ER Диаграмма (простая структура)

```text
users
  ↳ role → roles.name
roles
  ↕ (многие ко многим)
permissions
  ↕
role_permissions
```
