package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey")

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"` // можно оставить пустым
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	role := input.Role
	if role == "" {
		role = "user"
	}

	_, err := DB.Exec(`INSERT INTO users (username, password, role) VALUES ($1, $2, $3)`,
		input.Username, string(hashedPassword), role)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user User
	err := DB.Get(&user, `SELECT * FROM users WHERE username = $1`, input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create access token
	accessToken, err := CreateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token creation failed"})
		return
	}

	// Generate refresh token
	refreshToken, refreshExp, err := GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token creation failed"})
		return
	}

	// Вычисляем maxAge для cookie
	maxAge := int(time.Until(refreshExp).Seconds())

	// Кладём refresh в HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"",    // или твой домен в проде
		false, // secure: true на https
		true,  // HttpOnly
	)

	// Фронту отдаём access токен
	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

func ProtectedEndpoint(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{"message": "Welcome, " + username})
}

func AdminEndpoint(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{"message": "Hello, admin " + username})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Username        string `json:"username"`
		NewPassword     string `json:"new_password"`
		SuperAdminCode  string `json:"superadmin_code"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	expectedCode := os.Getenv("SUPERADMIN_CODE")
	if input.SuperAdminCode != expectedCode {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid SuperAdminCode"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	res, err := DB.Exec(`UPDATE users SET password=$1 WHERE username=$2`, string(hashedPassword), input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func MeEndpoint(c *gin.Context) {
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)

	// Получаем права доступа для роли пользователя
	var permissions []string
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN roles r ON rp.role_id = r.id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE r.name = $1
		ORDER BY p.name
	`
	err := DB.Select(&permissions, query, role)
	if err != nil {
		// Логируем ошибку для отладки
		log.Printf("Ошибка получения прав для роли %s: %v", role, err)
		permissions = []string{}
	}

	// Убеждаемся, что permissions не nil
	if permissions == nil {
		permissions = []string{}
	}

	// Возвращаем информацию о пользователе с правами
	c.JSON(http.StatusOK, gin.H{
		"username":    username,
		"role":        role,
		"permissions": permissions,
	})
}

func GenerateRefreshToken(userID int) (string, time.Time, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", time.Time{}, err
	}

	token := hex.EncodeToString(b)
	expires := time.Now().Add(30 * 24 * time.Hour)

	if err := SaveRefreshToken(userID, token, expires); err != nil {
		return "", time.Time{}, err
		}

	return token, expires, nil
}

func CreateAccessToken(user User) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // 15 минут жизни

	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func RefreshEndpoint(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	rt, err := GetRefreshToken(refreshToken)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	var user User
	if err := DB.Get(&user, "SELECT * FROM users WHERE id = $1", rt.UserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	accessToken, err := CreateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

func LogoutEndpoint(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		_ = DeleteRefreshToken(refreshToken)
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// Получение всех ролей и прав
func GetRoles(c *gin.Context) {
	var roles []Role
	err := DB.Select(&roles, "SELECT * FROM roles")
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	response := []gin.H{}
	for _, role := range roles {
		var perms []string
		query := `
			SELECT p.name
			FROM role_permissions rp
			JOIN permissions p ON rp.permission_id = p.id
			WHERE rp.role_id = $1
		`
		_ = DB.Select(&perms, query, role.ID)

		response = append(response, gin.H{
			"role":        role.Name,
			"description": role.Description,
			"permissions": perms,
		})
	}

	c.JSON(200, response)
}


func UpdateRolePermissions(c *gin.Context) {
	var input struct {
		SuperAdminCode string `json:"superadmin_code"`
		Updates []struct {
			Role        string   `json:"role"`
			Permissions []string `json:"permissions"`
		} `json:"updates"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if input.SuperAdminCode != os.Getenv("SUPERADMIN_CODE") {
		c.JSON(403, gin.H{"error": "Access denied"})
		return
	}

	for _, u := range input.Updates {
		var roleID int
		err := DB.Get(&roleID, "SELECT id FROM roles WHERE name=$1", u.Role)
		if err != nil {
			c.JSON(400, gin.H{"error": "Role not found: " + u.Role})
			return
		}

		// Очистить старые права
		DB.Exec("DELETE FROM role_permissions WHERE role_id=$1", roleID)

		// Назначить новые права
		for _, permName := range u.Permissions {
			var permID int
			err = DB.Get(&permID, "SELECT id FROM permissions WHERE name=$1", permName)
			if err != nil {
				c.JSON(400, gin.H{"error": "Permission not found: " + permName})
				return
			}
			DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)", roleID, permID)
		}
	}

	c.JSON(200, gin.H{"message": "Roles updated successfully"})
}
