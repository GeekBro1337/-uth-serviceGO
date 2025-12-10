package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalln("DB_DSN not set in .env")
	}

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("DB connection failed:", err)
	}

	// --- Миграции таблиц ---
	schemaUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);`

	schemaRoles := `
	CREATE TABLE IF NOT EXISTS roles (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		description TEXT
	);`

	schemaPermissions := `
	CREATE TABLE IF NOT EXISTS permissions (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		description TEXT
	);`

	schemaRolePermissions := `
	CREATE TABLE IF NOT EXISTS role_permissions (
		role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
		permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
		PRIMARY KEY (role_id, permission_id)
	);`

	schemaRefreshTokens := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL
	);`

	DB.MustExec(schemaUsers)
	DB.MustExec(schemaRoles)
	DB.MustExec(schemaPermissions)
	DB.MustExec(schemaRolePermissions)
	DB.MustExec(schemaRefreshTokens)

	fmt.Println("✅ DB connected and migrated")

	// --- Инициализация базовых ролей и прав ---
	var count int
	DB.Get(&count, "SELECT COUNT(*) FROM roles")
	if count == 0 {
		fmt.Println("⚙️ Initializing default roles and permissions...")

		// Добавляем роли
		DB.Exec(`INSERT INTO roles (name, description) VALUES
			('admin', 'Full access to everything'),
			('user', 'Standard user with limited rights'),
			('auditor', 'Read-only access')`)

		// Добавляем разрешения
		DB.Exec(`INSERT INTO permissions (name, description) VALUES
			('read', 'Read data'),
			('edit', 'Edit data'),
			('delete', 'Delete data'),
			('create', 'Create new records')`)

		// Привязываем права к ролям
		DB.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT r.id, p.id FROM roles r, permissions p
			WHERE r.name='admin';`) // admin получает всё

		DB.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT r.id, p.id FROM roles r, permissions p
			WHERE r.name='user' AND p.name = 'read';`)

		DB.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT r.id, p.id FROM roles r, permissions p
			WHERE r.name='auditor' AND p.name='read';`)

		fmt.Println("✅ Default roles and permissions initialized.")
	}
}

func SaveRefreshToken(userID int, token string, expires time.Time) error {
	_, err := DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expires)
	return err
}

func GetRefreshToken(token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := DB.Get(&rt, `
		SELECT * FROM refresh_tokens WHERE token = $1
	`, token)
	return &rt, err
}

func DeleteRefreshToken(token string) error {
	_, err := DB.Exec(`
		DELETE FROM refresh_tokens WHERE token = $1
	`, token)
	return err
}
