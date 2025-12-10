package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sqlx.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalln("DB_DSN not set in .env")
	}
	masterDSN := os.Getenv("DB_DSN_MASTER") // опционально: строка подключения к master

	// Вытаскиваем имя базы из DSN
	dbName := extractDBName(dsn)

	var err error
	DB, err = sqlx.Connect("sqlserver", dsn)
	if err != nil && masterDSN != "" && dbName != "" {
		// Если не удалось подключиться к целевой базе — пробуем создать её через master
		log.Printf("⚠️  DB connect failed: %v. Trying to create database %s via master...\n", err, dbName)
		if err := createDatabaseIfMissing(masterDSN, dbName); err != nil {
			log.Fatalln("DB create failed:", err)
		}
		DB, err = sqlx.Connect("sqlserver", dsn)
	}
	if err != nil {
		log.Fatalln("DB connection failed:", err)
	}

	// --- Миграции таблиц (SQL Server) ---
	schemaUsers := `
IF OBJECT_ID('users', 'U') IS NULL
BEGIN
    CREATE TABLE users (
        id INT IDENTITY(1,1) PRIMARY KEY,
        username NVARCHAR(255) NOT NULL UNIQUE,
        password NVARCHAR(255) NOT NULL,
        role NVARCHAR(50) NOT NULL DEFAULT 'user'
    );
END;`

	schemaRoles := `
IF OBJECT_ID('roles', 'U') IS NULL
BEGIN
    CREATE TABLE roles (
        id INT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(100) NOT NULL UNIQUE,
        description NVARCHAR(255)
    );
END;`

	schemaPermissions := `
IF OBJECT_ID('permissions', 'U') IS NULL
BEGIN
    CREATE TABLE permissions (
        id INT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(100) NOT NULL UNIQUE,
        description NVARCHAR(255)
    );
END;`

	schemaRolePermissions := `
IF OBJECT_ID('role_permissions', 'U') IS NULL
BEGIN
    CREATE TABLE role_permissions (
        role_id INT NOT NULL,
        permission_id INT NOT NULL,
        CONSTRAINT PK_role_permissions PRIMARY KEY (role_id, permission_id),
        CONSTRAINT FK_rp_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
        CONSTRAINT FK_rp_perm FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
    );
END;`

	schemaRefreshTokens := `
IF OBJECT_ID('refresh_tokens', 'U') IS NULL
BEGIN
    CREATE TABLE refresh_tokens (
        id INT IDENTITY(1,1) PRIMARY KEY,
        user_id INT NOT NULL,
        token NVARCHAR(512) NOT NULL UNIQUE,
        expires_at DATETIME2 NOT NULL,
        CONSTRAINT FK_rt_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
END;`

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

// extractDBName достаёт параметр database из DSN (sqlserver://...?...).
func extractDBName(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return ""
	}
	q := u.Query()
	return q.Get("database")
}

// createDatabaseIfMissing подключается к master и создаёт базу, если её нет.
func createDatabaseIfMissing(masterDSN, dbName string) error {
	masterDB, err := sqlx.Connect("sqlserver", masterDSN)
	if err != nil {
		return fmt.Errorf("connect master failed: %w", err)
	}
	defer masterDB.Close()

	_, err = masterDB.Exec(fmt.Sprintf("IF DB_ID('%s') IS NULL CREATE DATABASE %s;", dbName, dbName))
	return err
}

func SaveRefreshToken(userID int, token string, expires time.Time) error {
	_, err := DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES (@p1, @p2, @p3)
	`, userID, token, expires)
	return err
}

func GetRefreshToken(token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := DB.Get(&rt, `
		SELECT * FROM refresh_tokens WHERE token = @p1
	`, token)
	return &rt, err
}

func DeleteRefreshToken(token string) error {
	_, err := DB.Exec(`
		DELETE FROM refresh_tokens WHERE token = @p1
	`, token)
	return err
}
