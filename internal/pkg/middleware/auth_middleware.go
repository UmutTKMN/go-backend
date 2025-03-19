package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/UmutTKMN/go-backend/configs"
	"github.com/UmutTKMN/go-backend/internal/app/model"
	"github.com/UmutTKMN/go-backend/internal/app/services"
	"github.com/UmutTKMN/go-backend/internal/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT token'ı doğrular ve kullanıcı bilgilerini context'e ekler
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header'ından token'ı al
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token bulunamadı"})
			c.Abort()
			return
		}

		// Bearer şemasını kontrol et
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token formatı"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		config := configs.GetConfig()

		// Token'ı doğrula
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("geçersiz imza yöntemi")
			}
			return []byte(config.JWTKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token: " + err.Error()})
			c.Abort()
			return
		}

		// Token içeriğini al
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Kullanıcı ID'sini al
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token içeriği"})
				c.Abort()
				return
			}

			// Kullanıcı bilgilerini veritabanından al
			var user model.User
			if err := database.DB.First(&user, uint(userID)).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bulunamadı"})
				c.Abort()
				return
			}

			// Kullanıcı aktif mi kontrol et
			if !user.IsActive {
				c.JSON(http.StatusForbidden, gin.H{"error": "Hesabınız aktif değil"})
				c.Abort()
				return
			}

			// Kullanıcıyı context'e ekle
			c.Set("user", user)
			c.Set("userID", uint(userID))

			// Son aktivite zamanını güncelle
			database.DB.Model(&user).Update("last_activity", database.DB.NowFunc())

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			c.Abort()
			return
		}
	}
}

// GetCurrentUser context'ten mevcut kullanıcıyı alır
func GetCurrentUser(c *gin.Context) (model.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return model.User{}, false
	}
	return user.(model.User), true
}

// GetCurrentUserID context'ten mevcut kullanıcı ID'sini alır
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// RoleMiddleware yetki kontrolü için middleware
type RoleMiddleware struct {
	roleService *services.RoleService
}

// NewRoleMiddleware yeni bir RoleMiddleware örneği oluşturur
func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{
		roleService: services.NewRoleService(),
	}
}

// RequireAuth kullanıcının giriş yapmış olmasını gerektirir
// JWTAuth middleware'inden sonra çalıştırılmalıdır
func (m *RoleMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWTAuth tarafından ayarlanan userID kontrolü
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bu işlem için giriş yapmalısınız"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole kullanıcının belirtilen role sahip olmasını gerektirir
func (m *RoleMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWTAuth tarafından ayarlanan userID'yi al
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bu işlem için giriş yapmalısınız"})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı ID'si alınamadı"})
			c.Abort()
			return
		}

		// Rol kontrolü
		hasRole, err := m.roleService.HasRole(userID, roleName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol kontrolü yapılırken hata oluştu"})
			c.Abort()
			return
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu işlem için gereken role sahip değilsiniz: " + roleName})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin kullanıcının Admin rolüne sahip olmasını gerektirir
func (m *RoleMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("Admin")
}

// RequireSuperAdmin kullanıcının Super Admin rolüne sahip olmasını gerektirir
func (m *RoleMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return m.RequireRole("Super Admin")
}

// RequireStaff kullanıcının Staff rolüne sahip olmasını gerektirir
func (m *RoleMiddleware) RequireStaff() gin.HandlerFunc {
	return m.RequireRole("Staff")
}

// RequireManager kullanıcının Manager rolüne sahip olmasını gerektirir
func (m *RoleMiddleware) RequireManager() gin.HandlerFunc {
	return m.RequireRole("Manager")
}
