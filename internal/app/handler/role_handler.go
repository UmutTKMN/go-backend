package handler

import (
	"net/http"
	"strconv"

	"github.com/UmutTKMN/go-backend/internal/app/model"
	"github.com/UmutTKMN/go-backend/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RoleHandler rol işlemleri için handler
type RoleHandler struct {
	roleService *services.RoleService
}

// NewRoleHandler yeni bir RoleHandler örneği oluşturur
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

// GetAllRoles tüm rolleri getirir
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Roller getirilirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// GetRoleByID ID'ye göre bir rol getirir
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz rol ID'si"})
		return
	}

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// CreateRole yeni bir rol oluşturur
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	// Rol oluşturanı ayarla (isteğe bağlı)
	userID, exists := c.Get("userID")
	if exists {
		role.CreatedBy = userID.(uint)
	}

	if err := h.roleService.CreateRole(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol oluşturulurken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": role, "message": "Rol başarıyla oluşturuldu"})
}

// UpdateRole bir rolü günceller
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz rol ID'si"})
		return
	}

	// Önce mevcut rolü al
	existingRole, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Yeni verileri bağla
	var updatedRole model.Role
	if err := c.ShouldBindJSON(&updatedRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	// ID'yi ayarla ve güncelle
	updatedRole.ID = uint(id)

	// Sistem rolleri için kısıtlamalar
	if existingRole.IsSystemRole {
		// Sistem rolleri için isim değiştirmeyi önle
		updatedRole.RoleName = existingRole.RoleName
		// Sistem rolü özelliğini değiştirmeyi önle
		updatedRole.IsSystemRole = true
	}

	if err := h.roleService.UpdateRole(&updatedRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol güncellenirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedRole, "message": "Rol başarıyla güncellendi"})
}

// DeleteRole bir rolü siler
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz rol ID'si"})
		return
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol başarıyla silindi"})
}

// AssignRoleToUser kullanıcıya bir rol atar
func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	if err := h.roleService.AssignRoleToUser(request.UserID, request.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol atanırken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol başarıyla kullanıcıya atandı"})
}

// RemoveRoleFromUser kullanıcıdan bir rolü kaldırır
func (h *RoleHandler) RemoveRoleFromUser(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	if err := h.roleService.RemoveRoleFromUser(request.UserID, request.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol kaldırılırken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol başarıyla kullanıcıdan kaldırıldı"})
}

// GetUserRoles kullanıcının rollerini getirir
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si"})
		return
	}

	roles, err := h.roleService.GetUserRoles(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// CheckUserRole kullanıcının belirli bir role sahip olup olmadığını kontrol eder
func (h *RoleHandler) CheckUserRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si"})
		return
	}

	var request struct {
		RoleName string `json:"role_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	hasRole, err := h.roleService.HasRole(uint(userID), request.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol kontrolü yapılırken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_role": hasRole})
}
