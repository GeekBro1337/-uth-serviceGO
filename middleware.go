package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoggingMiddleware writes audit logs to DB (log_audit)
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()

		// берём username из контекста, если есть (после AuthMiddleware)
		username := c.GetString("username")

		// Определяем тип лога: авторизация или права (permissions)
		logType := "permissions"
		if strings.HasPrefix(path, "/login") ||
			strings.HasPrefix(path, "/register") ||
			strings.HasPrefix(path, "/refresh") ||
			strings.HasPrefix(path, "/logout") ||
			strings.HasPrefix(path, "/reset-password") {
			logType = "auth"
		}

		// сохраняем лог в отдельной горутине, чтобы не блокировать ответ
		go SaveAuditLog(path, method, status, int(latency.Milliseconds()), ip, username, logType)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func PermissionMiddleware(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName := c.GetString("role")

		var count int
		query := `
			SELECT COUNT(*)
			FROM role_permissions rp
			JOIN roles r ON rp.role_id = r.id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE r.name = @p1 AND p.name = @p2;
		`

		err := DB.Get(&count, query, roleName, requiredPerm)
		if err != nil || count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
