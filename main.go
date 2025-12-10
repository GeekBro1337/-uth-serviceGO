package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// Загружаем секрет для JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("⚠️  JWT_SECRET not set, using default secret (unsafe!)")
		jwtSecret = "supersecretkey"
	}
	jwtKey = []byte(jwtSecret)

	// Инициализация БД
	InitDB()

	// Инициализация роутера
	r := gin.Default()
	r.SetTrustedProxies(nil) // убираем предупреждение Gin о прокси

	// ✅ Аудит логов
	r.Use(LoggingMiddleware())

	// ✅ Разрешаем CORS для всех источников (с поддержкой credentials)
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // разрешаем все источники
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// r.Use(cors.New(cors.Config{
	// 	AllowOriginFunc: func(origin string) bool {
	// 		allowed := []string{
	// 			"http://dev.uktmp.kz:3002",
	// 			"http://dev.uktmp.kz:4040",
	// 			"http://localhost:3000",
	// 			"192.168.102.110",
	// 			"192.168.102.167",
	// 			"192.168.102.127",
	// 			"127.0.0.1",
	// 			"https://dev.uktmp.kz:4040",
	// 			"https://jwt.dev.uktmp.kz",
	// 		}
	// 		for _, a := range allowed {
	// 			if origin == a {
	// 				return true
	// 			}
	// 		}
	// 		return false
	// 	},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

// 	r.Use(cors.New(cors.Config{
//     AllowAllOrigins:  true, // ✅ разрешаем всем
//     AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
//     AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
//     AllowCredentials: true,
//     MaxAge:           12 * time.Hour,
// }))


	// --- Публичные маршруты ---
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.POST("/refresh", RefreshEndpoint)
	r.POST("/logout", LogoutEndpoint)
	r.POST("/reset-password", ResetPassword)

	// --- Защищённые маршруты ---
	protected := r.Group("/protected")
	protected.Use(AuthMiddleware())

	protected.GET("/user", RoleMiddleware("user"), ProtectedEndpoint)
	protected.GET("/admin", RoleMiddleware("admin"), AdminEndpoint)
	protected.GET("/me", MeEndpoint)

	// --- Проверка прав доступа (через PermissionMiddleware) ---
	protected.GET("/read", PermissionMiddleware("read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Read access granted"})
	})
	protected.GET("/edit", PermissionMiddleware("edit"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Edit access granted"})
	})
	protected.GET("/delete", PermissionMiddleware("delete"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "You can delete"})
	})

	// --- Управление ролями и правами ---
	r.GET("/roles", GetRoles)
	r.POST("/roles/update", UpdateRolePermissions) // добавлен слеш перед roles

	// --- Запуск сервера ---
	log.Println("🚀 Server started on :8999")
	r.Run(":8999")
}
