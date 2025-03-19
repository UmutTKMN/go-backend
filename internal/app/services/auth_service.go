package services

import (
	"errors"
	"time"

	"github.com/UmutTKMN/go-backend/internal/app/model"
	"github.com/UmutTKMN/go-backend/internal/pkg/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	secretKey []byte
}

func NewAuthService(secretKey string) *AuthService {
	return &AuthService{
		secretKey: []byte(secretKey),
	}
}

// Register yeni kullanıcı kaydı yapar
func (s *AuthService) Register(req *model.RegisterRequest) (*model.User, error) {
	// Email veya kullanıcı adı zaten var mı kontrol et
	var existingUser model.User
	result := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		if existingUser.Email == req.Email {
			return nil, errors.New("bu email adresi zaten kullanılıyor")
		}
		return nil, errors.New("bu kullanıcı adı zaten kullanılıyor")
	}

	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Doğum tarihi işleme
	var birthDate *time.Time
	if req.BirthDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			birthDate = &parsedTime
		}
	}

	now := time.Now()
	// Yeni kullanıcı oluştur
	user := &model.User{
		Username:         req.Username,
		Email:            req.Email,
		Password:         string(hashedPassword),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		PhoneNumber:      req.PhoneNumber,
		Country:          req.Country,
		BirthDate:        birthDate,
		CreatedAt:        now,
		UpdatedAt:        now,
		RegistrationDate: now,
		IsActive:         true,
		IsVerified:       false,
		// Benzersiz doğrulama tokeni oluştur
		VerificationToken:  generateRandomToken(),
		VerificationSentAt: &now,
		PreferredLanguage:  "tr",
		Timezone:           "Europe/Istanbul",
		ThemePreference:    "light",
		AccountType:        "standard",
	}

	// Görünen adı otomatik ayarla
	if user.DisplayName == "" {
		if user.FirstName != "" && user.LastName != "" {
			user.DisplayName = user.FirstName + " " + user.LastName
		} else {
			user.DisplayName = user.Username
		}
	}

	// Veritabanına kaydet
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Doğrulama e-postası gönder
	// TODO: Implement email sending

	return user, nil
}

// Login kullanıcı girişi yapar ve token döndürür
func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	// Veritabanından kullanıcıyı bul
	var user model.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	// Kullanıcı aktif mi kontrol et
	if !user.IsActive {
		return nil, errors.New("hesabınız aktif değil")
	}

	// Hesap kilitli mi kontrol et
	if user.AccountLockedUntil != nil && user.AccountLockedUntil.After(time.Now()) {
		return nil, errors.New("hesabınız geçici olarak kilitlendi, lütfen daha sonra tekrar deneyin")
	}

	// Şifre kontrolü
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Başarısız giriş denemesini kaydet
		user.LoginAttempts++

		// 5 başarısız denemeden sonra hesabı geçici olarak kilitle
		if user.LoginAttempts >= 5 {
			lockUntil := time.Now().Add(time.Minute * 30) // 30 dakika kilitle
			user.AccountLockedUntil = &lockUntil
		}

		database.DB.Save(&user)
		return nil, errors.New("geçersiz şifre")
	}

	// Başarılı giriş işlemleri
	now := time.Now()
	user.LastLogin = &now
	user.LastActivity = &now
	user.LoginAttempts = 0        // Başarılı giriş, sayacı sıfırla
	user.AccountLockedUntil = nil // Kilidi kaldır

	// Kullanıcıyı güncelle
	database.DB.Save(&user)

	// JWT token oluştur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

// Rastgele doğrulama tokeni oluşturmak için yardımcı fonksiyon
func generateRandomToken() string {
	// Gerçek uygulamada güvenli bir rastgele dize oluşturma algoritması kullanılmalı
	return "verification-token-" + time.Now().Format("20060102150405")
}
