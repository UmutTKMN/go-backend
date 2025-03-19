package services

import (
	"errors"
	"time"

	"github.com/UmutTKMN/go-backend/internal/app/model"
	"github.com/UmutTKMN/go-backend/internal/pkg/database"
	"gorm.io/gorm"
)

// StaffService personel işlemleri için servis
type StaffService struct {
	db *gorm.DB
}

// NewStaffService yeni bir StaffService örneği oluşturur
func NewStaffService() *StaffService {
	return &StaffService{
		db: database.DB,
	}
}

// GetAllStaff tüm personel bilgilerini getirir
func (s *StaffService) GetAllStaff() ([]model.Staff, error) {
	var staffList []model.Staff

	// İlişkili verileri de yükle
	result := s.db.Preload("User").Preload("Role").Preload("Manager").Find(&staffList)
	return staffList, result.Error
}

// GetStaffByID ID'ye göre personel bilgisi getirir
func (s *StaffService) GetStaffByID(id uint) (*model.Staff, error) {
	var staff model.Staff

	// İlişkili verileri de yükle
	result := s.db.Preload("User").Preload("Role").Preload("Manager").First(&staff, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("personel bulunamadı")
		}
		return nil, result.Error
	}

	return &staff, nil
}

// GetStaffByUserID kullanıcı ID'sine göre personel bilgisi getirir
func (s *StaffService) GetStaffByUserID(userID uint) (*model.Staff, error) {
	var staff model.Staff

	// İlişkili verileri de yükle
	result := s.db.Preload("User").Preload("Role").Preload("Manager").Where("user_id = ?", userID).First(&staff)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcıya ait personel kaydı bulunamadı")
		}
		return nil, result.Error
	}

	return &staff, nil
}

// CreateStaff yeni bir personel kaydı oluşturur
func (s *StaffService) CreateStaff(staff *model.Staff) error {
	// Kullanıcının var olup olmadığını kontrol et
	var user model.User
	if err := s.db.First(&user, staff.UserID).Error; err != nil {
		return errors.New("kullanıcı bulunamadı")
	}

	// Rolün var olup olmadığını kontrol et
	var role model.Role
	if err := s.db.First(&role, staff.RoleID).Error; err != nil {
		return errors.New("rol bulunamadı")
	}

	// Yönetici belirtilmişse var olup olmadığını kontrol et
	if staff.ManagerID != nil && *staff.ManagerID > 0 {
		var manager model.Staff
		if err := s.db.First(&manager, *staff.ManagerID).Error; err != nil {
			return errors.New("yönetici bulunamadı")
		}
	}

	// Başlangıç tarihini ayarla
	if staff.StartDate.IsZero() {
		staff.StartDate = time.Now()
	}

	// Kullanıcıya rolü ekle (many-to-many ilişki)
	if err := s.db.Model(&user).Association("Roles").Append(&role); err != nil {
		return err
	}

	// Personel kaydını oluştur
	return s.db.Create(staff).Error
}

// UpdateStaff bir personel kaydını günceller
func (s *StaffService) UpdateStaff(staff *model.Staff) error {
	// Personelin var olup olmadığını kontrol et
	result := s.db.First(&model.Staff{}, staff.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("personel bulunamadı")
		}
		return result.Error
	}

	// Rol değişmişse ve geçerliyse
	if staff.RoleID > 0 {
		var role model.Role
		if err := s.db.First(&role, staff.RoleID).Error; err != nil {
			return errors.New("rol bulunamadı")
		}

		// Kullanıcının rollerini güncelle
		var user model.User
		if err := s.db.First(&user, staff.UserID).Error; err != nil {
			return errors.New("kullanıcı bulunamadı")
		}

		// Eski rolü kaldır ve yeni rolü ekle
		// NOT: Bu örnekte basitleştirmek için tüm rolleri kaldırıp yeni rolü ekliyoruz
		// Gerçek uygulamada daha karmaşık bir rol yönetimi gerekebilir
		if err := s.db.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}

		if err := s.db.Model(&user).Association("Roles").Append(&role); err != nil {
			return err
		}
	}

	// Yönetici belirtilmişse var olup olmadığını kontrol et
	if staff.ManagerID != nil && *staff.ManagerID > 0 {
		var manager model.Staff
		if err := s.db.First(&manager, *staff.ManagerID).Error; err != nil {
			return errors.New("yönetici bulunamadı")
		}
	}

	// Personel kaydını güncelle
	return s.db.Save(staff).Error
}

// DeleteStaff bir personel kaydını siler
func (s *StaffService) DeleteStaff(id uint) error {
	// Personel var mı kontrol et
	var staff model.Staff
	if err := s.db.First(&staff, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("personel bulunamadı")
		}
		return err
	}

	// Kullanıcının rollerini temizle
	var user model.User
	if err := s.db.First(&user, staff.UserID).Error; err == nil {
		if err := s.db.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
	}

	// Personeli sil
	return s.db.Delete(&staff).Error
}

// GetStaffByDepartment departmana göre personel listesi getirir
func (s *StaffService) GetStaffByDepartment(department string) ([]model.Staff, error) {
	var staffList []model.Staff

	result := s.db.Preload("User").Preload("Role").Preload("Manager").
		Where("department = ?", department).
		Find(&staffList)

	return staffList, result.Error
}

// GetStaffByManager yöneticiye göre personel listesi getirir
func (s *StaffService) GetStaffByManager(managerID uint) ([]model.Staff, error) {
	var staffList []model.Staff

	result := s.db.Preload("User").Preload("Role").Preload("Manager").
		Where("manager_id = ?", managerID).
		Find(&staffList)

	return staffList, result.Error
}

// GetStaffByRole role göre personel listesi getirir
func (s *StaffService) GetStaffByRole(roleID uint) ([]model.Staff, error) {
	var staffList []model.Staff

	result := s.db.Preload("User").Preload("Role").Preload("Manager").
		Where("role_id = ?", roleID).
		Find(&staffList)

	return staffList, result.Error
}

// UpdateStaffPosition personelin pozisyonunu günceller
func (s *StaffService) UpdateStaffPosition(staffID uint, department, position string) error {
	var staff model.Staff
	if err := s.db.First(&staff, staffID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("personel bulunamadı")
		}
		return err
	}

	// Sadece belirli alanları güncelle
	return s.db.Model(&staff).Updates(map[string]interface{}{
		"department": department,
		"position":   position,
	}).Error
}

// UpdateStaffManager personelin yöneticisini günceller
func (s *StaffService) UpdateStaffManager(staffID, managerID uint) error {
	var staff model.Staff
	if err := s.db.First(&staff, staffID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("personel bulunamadı")
		}
		return err
	}

	if managerID > 0 {
		var manager model.Staff
		if err := s.db.First(&manager, managerID).Error; err != nil {
			return errors.New("yönetici bulunamadı")
		}

		// Kendisini yönetici olarak atama kontrolü
		if staffID == managerID {
			return errors.New("personel kendisini yönetici olarak atayamaz")
		}
	}

	return s.db.Model(&staff).Update("manager_id", managerID).Error
}
