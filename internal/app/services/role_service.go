package services

import (
	"errors"

	"github.com/UmutTKMN/gobackend/internal/app/model"
	"github.com/UmutTKMN/gobackend/internal/pkg/database"
	"gorm.io/gorm"
)

// RoleService rol işlemleri için servis
type RoleService struct {
	db *gorm.DB
}

// NewRoleService yeni bir RoleService örneği oluşturur
func NewRoleService() *RoleService {
	return &RoleService{
		db: database.DB,
	}
}

// GetAllRoles tüm rolleri getirir
func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	result := s.db.Find(&roles)
	return roles, result.Error
}

// GetRoleByID ID'ye göre bir rol getirir
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	result := s.db.First(&role, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("rol bulunamadı")
		}
		return nil, result.Error
	}
	return &role, nil
}

// CreateRole yeni bir rol oluşturur
func (s *RoleService) CreateRole(role *model.Role) error {
	return s.db.Create(role).Error
}

// UpdateRole bir rolü günceller
func (s *RoleService) UpdateRole(role *model.Role) error {
	result := s.db.First(&model.Role{}, role.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("rol bulunamadı")
		}
		return result.Error
	}

	return s.db.Save(role).Error
}

// DeleteRole bir rolü siler
func (s *RoleService) DeleteRole(id uint) error {
	// Sistem rollerini silmeyi engelle
	var role model.Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rol bulunamadı")
		}
		return err
	}

	if role.IsSystemRole {
		return errors.New("sistem rolleri silinemez")
	}

	// Rol kullanılıyor mu kontrol et
	var count int64
	if err := s.db.Model(&model.Staff{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("bu rol hala personel tarafından kullanılıyor ve silinemez")
	}

	return s.db.Delete(&model.Role{}, id).Error
}

// AssignRoleToUser kullanıcıya bir rol atar
func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("kullanıcı bulunamadı")
	}

	var role model.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return errors.New("rol bulunamadı")
	}

	// Many-to-many ilişki
	return s.db.Model(&user).Association("Roles").Append(&role)
}

// RemoveRoleFromUser kullanıcıdan bir rolü kaldırır
func (s *RoleService) RemoveRoleFromUser(userID, roleID uint) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("kullanıcı bulunamadı")
	}

	var role model.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return errors.New("rol bulunamadı")
	}

	return s.db.Model(&user).Association("Roles").Delete(&role)
}

// GetUserRoles kullanıcının rollerini getirir
func (s *RoleService) GetUserRoles(userID uint) ([]model.Role, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	var roles []model.Role
	if err := s.db.Model(&user).Association("Roles").Find(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// HasRole kullanıcının belirli bir role sahip olup olmadığını kontrol eder
func (s *RoleService) HasRole(userID uint, roleName string) (bool, error) {
	var count int64
	err := s.db.Raw(`
		SELECT COUNT(*) FROM users u 
		JOIN user_roles ur ON u.id = ur.user_id 
		JOIN roles r ON ur.role_id = r.role_id 
		WHERE u.id = ? AND r.role_name = ?
	`, userID, roleName).Count(&count).Error

	return count > 0, err
}

// UserRoleDetail kullanıcı-rol ilişkisi için detaylı veri yapısı
type UserRoleDetail struct {
	UserID    uint   `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	RoleID    uint   `json:"role_id"`
	RoleName  string `json:"role_name"`
	Level     int    `json:"permission_level"`
}

// GetAllUserRoleDetails tüm kullanıcı-rol ilişkilerini detaylı olarak getirir
func (s *RoleService) GetAllUserRoleDetails() ([]UserRoleDetail, error) {
	var details []UserRoleDetail
	err := s.db.Raw(`
		SELECT 
			u.id as user_id, 
			u.email as user_email, 
			u.name as user_name, 
			r.role_id as role_id, 
			r.role_name as role_name,
			r.permission_level as level
		FROM 
			users u
		JOIN 
			user_roles ur ON u.id = ur.user_id
		JOIN 
			roles r ON ur.role_id = r.role_id
		ORDER BY 
			u.id, r.permission_level DESC
	`).Scan(&details).Error

	return details, err
}

// GetUserRoleDetailsByUser belirli bir kullanıcının rol detaylarını getirir
func (s *RoleService) GetUserRoleDetailsByUser(userID uint) ([]UserRoleDetail, error) {
	var details []UserRoleDetail
	err := s.db.Raw(`
		SELECT 
			u.id as user_id, 
			u.email as user_email, 
			u.name as user_name, 
			r.role_id as role_id, 
			r.role_name as role_name,
			r.permission_level as level
		FROM 
			users u
		JOIN 
			user_roles ur ON u.id = ur.user_id
		JOIN 
			roles r ON ur.role_id = r.role_id
		WHERE 
			u.id = ?
		ORDER BY 
			r.permission_level DESC
	`, userID).Scan(&details).Error

	return details, err
}

// GetUserRoleDetailsByRole belirli bir role sahip kullanıcıların detaylarını getirir
func (s *RoleService) GetUserRoleDetailsByRole(roleID uint) ([]UserRoleDetail, error) {
	var details []UserRoleDetail
	err := s.db.Raw(`
		SELECT 
			u.id as user_id, 
			u.email as user_email, 
			u.name as user_name, 
			r.role_id as role_id, 
			r.role_name as role_name,
			r.permission_level as level
		FROM 
			users u
		JOIN 
			user_roles ur ON u.id = ur.user_id
		JOIN 
			roles r ON ur.role_id = r.role_id
		WHERE 
			r.role_id = ?
		ORDER BY 
			u.id
	`, roleID).Scan(&details).Error

	return details, err
}
