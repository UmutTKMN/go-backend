package main

import (
	"log"
	"net/http"

	"github.com/UmutTKMN/go-backend/configs"
	"github.com/UmutTKMN/go-backend/internal/app/handler"
	"github.com/UmutTKMN/go-backend/internal/app/services"
	"github.com/UmutTKMN/go-backend/internal/pkg/database"
	"github.com/UmutTKMN/go-backend/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Yapılandırmayı yükle
	config := configs.GetConfig()

	// Veritabanını başlat
	database.Init(config)

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// Güvenilir proxy'leri ayarla
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Servisleri oluştur
	authService := services.NewAuthService(config.JWTKey)
	userService := services.NewUserService()

	// Handler'ları oluştur
	userHandler := handler.NewUserHandler(userService, authService)
	profileHandler := handler.NewProfileHandler(userService)

	// API v1 grubu
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hoş geldiniz!",
				"version": "1.0",
				"status":  "active",
			})
		})

		// Auth endpoint'leri
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Kullanıcı endpoint'leri
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Profil yönetimi rotaları
		profileGroup := v1.Group("/profile")
		profileGroup.Use(middleware.AuthMiddleware())
		{
			profileGroup.GET("", profileHandler.GetProfile)
			profileGroup.PUT("", profileHandler.UpdateProfile)
		}

		// Rol ve Personel middleware tanımla
		roleMiddleware := middleware.NewRoleMiddleware()

		// Rol yönetimi rotaları
		roleHandler := handler.NewRoleHandler()
		roleGroup := v1.Group("/roles")
		roleGroup.Use(middleware.AuthMiddleware())
		{
			roleGroup.GET("", roleHandler.GetAllRoles)

			// Admin gerektiren rotalar
			adminRoleGroup := roleGroup.Group("")
			adminRoleGroup.Use(roleMiddleware.RequireSuperAdmin())
			{
				adminRoleGroup.GET("/:id", roleHandler.GetRoleByID)
				adminRoleGroup.POST("", roleHandler.CreateRole)
				adminRoleGroup.POST("/assign", roleHandler.AssignRoleToUser)
				adminRoleGroup.POST("/remove", roleHandler.RemoveRoleFromUser)
				adminRoleGroup.PUT("/:id", roleHandler.UpdateRole)
				adminRoleGroup.DELETE("/:id", roleHandler.DeleteRole)
			}

			// Kullanıcı rolleri
			roleGroup.GET("/user/:id", roleHandler.GetUserRoles)
			roleGroup.POST("/user/:id/check-role", roleHandler.CheckUserRole)
		}

		// Personel yönetimi rotaları
		staffHandler := handler.NewStaffHandler()
		staffGroup := v1.Group("/staff")
		staffGroup.Use(middleware.AuthMiddleware())
		{
			// Yönetici gerektiren rotalar
			managerStaffGroup := staffGroup.Group("")
			managerStaffGroup.Use(roleMiddleware.RequireManager())
			{
				managerStaffGroup.GET("", staffHandler.GetAllStaff)
				managerStaffGroup.GET("/:id", staffHandler.GetStaffByID)
				managerStaffGroup.GET("/user/:id", staffHandler.GetStaffByUserID)
				managerStaffGroup.GET("/department/:department", staffHandler.GetStaffByDepartment)
				managerStaffGroup.GET("/manager/:id", staffHandler.GetStaffByManager)
				managerStaffGroup.GET("/role/:id", staffHandler.GetStaffByRole)
			}

			// Admin gerektiren rotalar
			adminStaffGroup := staffGroup.Group("")
			adminStaffGroup.Use(roleMiddleware.RequireAdmin())
			{
				adminStaffGroup.POST("", staffHandler.CreateStaff)
				adminStaffGroup.PUT("/:id", staffHandler.UpdateStaff)
				adminStaffGroup.PUT("/:id/position", staffHandler.UpdateStaffPosition)
				adminStaffGroup.PUT("/:id/manager", staffHandler.UpdateStaffManager)
				adminStaffGroup.DELETE("/:id", staffHandler.DeleteStaff)
			}
		}
	}

	// Sunucuyu başlat
	if err := router.Run(":" + config.AppPort); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
