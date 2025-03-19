package handler

import (
	"net/http"
	"strconv"

	"github.com/UmutTKMN/gobackend/internal/app/model"
	"github.com/UmutTKMN/gobackend/internal/app/services"
	"github.com/gin-gonic/gin"
)

// StaffHandler personel işlemleri için handler
type StaffHandler struct {
	staffService *services.StaffService
	userService  *services.UserService
}

// NewStaffHandler yeni bir StaffHandler örneği oluşturur
func NewStaffHandler() *StaffHandler {
	return &StaffHandler{
		staffService: services.NewStaffService(),
		userService:  services.NewUserService(),
	}
}

// GetAllStaff tüm personel bilgilerini getirir
func (h *StaffHandler) GetAllStaff(c *gin.Context) {
	staffList, err := h.staffService.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personel bilgileri getirilirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffList})
}

// GetStaffByID ID'ye göre personel bilgisi getirir
func (h *StaffHandler) GetStaffByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz personel ID'si"})
		return
	}

	staff, err := h.staffService.GetStaffByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staff})
}

// GetStaffByUserID kullanıcı ID'sine göre personel bilgisi getirir
func (h *StaffHandler) GetStaffByUserID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si"})
		return
	}

	staff, err := h.staffService.GetStaffByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staff})
}

// CreateStaff yeni bir personel kaydı oluşturur
func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var staff model.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	// Önce kullanıcının var olup olmadığını kontrol et
	_, err := h.userService.GetUserByID(staff.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı bulunamadı: " + err.Error()})
		return
	}

	if err := h.staffService.CreateStaff(&staff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personel kaydı oluşturulurken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": staff, "message": "Personel kaydı başarıyla oluşturuldu"})
}

// UpdateStaff bir personel kaydını günceller
func (h *StaffHandler) UpdateStaff(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz personel ID'si"})
		return
	}

	// Önce mevcut personeli al
	existingStaff, err := h.staffService.GetStaffByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Yeni verileri bağla
	var updatedStaff model.Staff
	if err := c.ShouldBindJSON(&updatedStaff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	// ID'yi ve UserID'yi değiştirmeye izin verme
	updatedStaff.ID = uint(id)
	updatedStaff.UserID = existingStaff.UserID

	if err := h.staffService.UpdateStaff(&updatedStaff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personel kaydı güncellenirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedStaff, "message": "Personel kaydı başarıyla güncellendi"})
}

// DeleteStaff bir personel kaydını siler
func (h *StaffHandler) DeleteStaff(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz personel ID'si"})
		return
	}

	if err := h.staffService.DeleteStaff(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personel kaydı başarıyla silindi"})
}

// GetStaffByDepartment departmana göre personel listesi getirir
func (h *StaffHandler) GetStaffByDepartment(c *gin.Context) {
	department := c.Param("department")
	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Departman belirtilmedi"})
		return
	}

	staffList, err := h.staffService.GetStaffByDepartment(department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personel bilgileri getirilirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffList})
}

// GetStaffByManager yöneticiye göre personel listesi getirir
func (h *StaffHandler) GetStaffByManager(c *gin.Context) {
	managerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz yönetici ID'si"})
		return
	}

	staffList, err := h.staffService.GetStaffByManager(uint(managerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personel bilgileri getirilirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffList})
}

// GetStaffByRole role göre personel listesi getirir
func (h *StaffHandler) GetStaffByRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz rol ID'si"})
		return
	}

	staffList, err := h.staffService.GetStaffByRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personel bilgileri getirilirken hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffList})
}

// UpdateStaffPosition personelin pozisyonunu günceller
func (h *StaffHandler) UpdateStaffPosition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz personel ID'si"})
		return
	}

	var request struct {
		Department string `json:"department" binding:"required"`
		Position   string `json:"position" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	if err := h.staffService.UpdateStaffPosition(uint(id), request.Department, request.Position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personel pozisyonu başarıyla güncellendi"})
}

// UpdateStaffManager personelin yöneticisini günceller
func (h *StaffHandler) UpdateStaffManager(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz personel ID'si"})
		return
	}

	var request struct {
		ManagerID uint `json:"manager_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri: " + err.Error()})
		return
	}

	if err := h.staffService.UpdateStaffManager(uint(id), request.ManagerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personel yöneticisi başarıyla güncellendi"})
}
