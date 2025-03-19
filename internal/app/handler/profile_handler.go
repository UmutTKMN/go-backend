package handler

import (
	"net/http"

	"github.com/UmutTKMN/go-backend/internal/app/services"
	"github.com/UmutTKMN/go-backend/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userService *services.UserService
}

func NewProfileHandler(userService *services.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// GetProfile kullanıcının kendi profilini görüntüler
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Context'ten kullanıcı bilgilerini al
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Kullanıcı servisinden güncel bilgileri al
	profile, err := h.userService.GetUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profil bilgileri alınamadı"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile kullanıcının kendi profilini günceller
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Context'ten kullanıcı ID'sini al
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Request body'den güncellenecek bilgileri al
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kullanıcı bilgilerini güncelle
	updatedUser, err := h.userService.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profil güncellenemedi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
