package services

import (
	"errors"
	"time"

	"github.com/UmutTKMN/gobackend/internal/app/model"
	"github.com/UmutTKMN/gobackend/internal/pkg/database"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// GetUser belirtilen ID'ye sahip kullanıcıyı döndürür
func (s *UserService) GetUser(id uint) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	// Yaşı hesapla
	if user.BirthDate != nil {
		now := time.Now()
		user.Age = now.Year() - user.BirthDate.Year()
		// Doğum günü bu yıl henüz gelmediyse yaşı bir azalt
		if now.YearDay() < user.BirthDate.YearDay() {
			user.Age--
		}
	}

	return &user, nil
}

// UpdateUser kullanıcı bilgilerini günceller
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*model.User, error) {
	// Önce kullanıcıyı bul
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	// Güvenli olmayan alanları temizle
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "created_at")
	delete(updates, "email_verified_at")
	delete(updates, "is_verified")
	delete(updates, "verification_token")
	delete(updates, "api_key")

	// Güncellemeleri uygula
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Güncellenmiş kullanıcıyı döndür
	return s.GetUser(id)
}

// DeleteUser kullanıcıyı siler
func (s *UserService) DeleteUser(id uint) error {
	result := database.DB.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("kullanıcı bulunamadı")
	}
	return nil
}

// GetAllUsers tüm kullanıcıları döndürür
func (s *UserService) GetAllUsers(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Toplam kayıt sayısını al
	if err := database.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sayfalama ile kullanıcıları al
	offset := (page - 1) * limit
	if err := database.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID ID'ye göre kullanıcı getirir
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, result.Error
	}
	return &user, nil
}
